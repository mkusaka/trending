package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Client struct {
	client *http.Client

	Logger *log.Logger
}

func NewClient(timeoutSecond time.Duration) *Client {
	client := &http.Client{
		Timeout: timeoutSecond,
	}
	return &Client{
		client: client,
		Logger: log.New(ioutil.Discard, "go-client: ", log.LstdFlags),
	}
}

type User struct {
	Username string
	Href     string
	Avatar   string
}

type Trend struct {
	Author             string
	Name               string
	Avatar             string
	Url                string
	Description        string
	Language           string
	LanguageColor      string
	Stars              int
	Forks              int
	CurrentPeriodStars int
	BuiltBy            []User
}

func ParseTrend(RawTrending string) []Trend {
	var Trends []Trend
	json.Unmarshal([]byte(RawTrending), &Trends)
	return Trends
}

func GenerateMarkDown(trends []Trend, trendTitle string) string {
	trendMD := "# " + trendTitle + "\n"
	for _, trend := range trends {
		if trendTitle == "general" {
			trendMD += "- [" + trend.Name + "](" + trend.Url + ") : " + trend.Language + "\n"
			if trend.Description != "" {
				trendMD += "  - " + trend.Description + "\n"
			}
		} else {
			trendMD += "- [" + trend.Name + "](" + trend.Url + ")\n  - " + trend.Description + "\n"
		}
	}
	return trendMD
}

var (
	General = ""

	Languages = []string{
		"go",
		"javascript",
		"ruby",
		"rust",
		"c++",
		"typescript",
	}

	Periods = []string{
		"daily",
		"weekly",
		"monthly",
	}
)

const baseUrl = "https://github-trending-api.now.sh/repositories"

// https://githubtrendingapi.docs.apiary.io/#reference/0/languages-collection/list-trending-repositories
func URLGenerator(period string) func(language string) string {
	return func(language string) string {
		return baseUrl + "?language=" + language + "&since=" + period
	}
}

func FetchAndGenerateMarkdown(period string) {
	c := NewClient(10 * time.Second)
	targets := append([]string{General}, Languages...)

	periodURLGenerator := URLGenerator(period)
	for _, language := range targets {
		url := periodURLGenerator(language)
		fmt.Print(url + "\n")
		resp, err := c.client.Get(url)
		if err != nil {
			log.Fatal("something error")
		}

		if resp.Body != nil {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Body io convert error")
			}
			parsed := ParseTrend(string(body))
			if language == "" {
				language = "general"
			}
			md := GenerateMarkDown(parsed, language)
			filename := "src/languages/" + language + "/" + period + ".md"
			f, err := os.Create(filename)
			if err != nil {
				log.Fatalf("something wrong with create md file %s", filename)
			}
			defer f.Close()

			f.WriteString(md)
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	period := os.Args[1]
	isValidPeriod := false
	for _, _period := range Periods {
		if period == _period {
			isValidPeriod = true
			break
		}
	}
	isAllPeriod := false
	if period == "all" {
		isAllPeriod = true
	}
	if !(isValidPeriod || isAllPeriod) {
		log.Fatal("invalid Period given. Choose valid period from daily, weekly, monthly or all.")
	}
	if isAllPeriod {
		for _, _period := range Periods {
			FetchAndGenerateMarkdown(_period)
		}
	} else {
		FetchAndGenerateMarkdown(period)
	}
}
