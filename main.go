package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

func GenerateMarkDown(trends []Trend, trendTitle string, depth int) string {
	trendMD := strings.Repeat("#", depth) + " " + trendTitle + "\n"
	for _, trend := range trends {
		if trendTitle == "general" {
			trendMD += "- [" + trend.Name + "](" + trend.Url + ") : " + trend.Language + "\n"
			if strings.TrimSpace(trend.Description) != "" {
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
		"typescript",
		"kotlin",
		"ruby",
		"rust",
		"c++",
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

func FetchAndGenerateJSON(period string) {
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
			if language == "" {
				language = "general"
			}
			jsonFilename := "src/raw/" + language + "/" + period + ".json"
			jsonfile, err := os.Create(jsonFilename)
			if err != nil {
				log.Fatalf("something wrong with create md file %s", jsonFilename)
			}
			defer jsonfile.Close()

			jsonfile.WriteString(string(body))
		}
		time.Sleep(1 * time.Second)
	}
}

func storeToLanguageMarkdown(language, period string) {
	jsonFilename := "src/raw/" + language + "/" + period + ".json"
	jsonByteString, err := ioutil.ReadFile(jsonFilename)
	if err != nil {
		log.Fatalf("something wrong with read file %s", jsonFilename)
	}

	parsed := ParseTrend(string(jsonByteString))
	md := GenerateMarkDown(parsed, language, 1)
	filename := "src/languages/" + language + "/" + period + ".md"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("something wrong with create md file %s", filename)
	}
	defer f.Close()

	f.WriteString(md)
}

func storeToPeriodMarkdown(period string) {
	filename := "src/periods/" + period + ".md"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("something wrong with create md file %s", filename)
	}
	defer f.Close()
	md := "# " + period + "\n"
	for _, language := range append([]string{General}, Languages...) {
		if language == "" {
			language = "general"
		}
		jsonFilename := "src/raw/" + language + "/" + period + ".json"
		jsonByteString, err := ioutil.ReadFile(jsonFilename)
		if err != nil {
			log.Fatalf("something wrong with read file %s", jsonFilename)
		}

		parsed := ParseTrend(string(jsonByteString))
		md += GenerateMarkDown(parsed, language, 2)
	}
	f.WriteString(md)
}

func FetchJSON(period string, isAllPeriod bool) {
	if isAllPeriod {
		for _, _period := range Periods {
			FetchAndGenerateJSON(_period)
		}
	} else {
		FetchAndGenerateJSON(period)
	}
}

func GeneratePeriodMarkdown(period string, isAllPeriod bool) {
	if isAllPeriod {
		for _, _period := range Periods {
			storeToPeriodMarkdown(_period)
		}
	} else {
		storeToPeriodMarkdown(period)
	}
}

func GenerateLanguageMarkdown(language string, isAllLanguages bool) {
	if isAllLanguages {
		for _, _language := range append([]string{General}, Languages...) {
			if _language == "" {
				_language = "general"
			}
			for _, period := range Periods {
				storeToLanguageMarkdown(_language, period)
			}
		}
	} else {
		for _, period := range Periods {
			storeToLanguageMarkdown(language, period)
		}
	}
}

func main() {
	period := os.Args[1]
	language := ""
	if len(os.Args) >= 3 {
		language = os.Args[2]
	}
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

	FetchJSON(period, isAllPeriod)
	GeneratePeriodMarkdown(period, isAllPeriod)

	if language != "" {
		isValidLanguage := false
		for _, _language := range Languages {
			if language == _language {
				isValidLanguage = true
				break
			}
		}
		if !isValidLanguage {
			log.Fatal("invalid language")
		}
		GenerateLanguageMarkdown(language, false)
	} else {
		GenerateLanguageMarkdown("", true)
	}
}
