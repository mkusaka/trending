package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Trend struct {
	Author      string
	Name        string
	Avatar      string
	Href        string
	Description string
	Language    string
	Stars       int
	Forks       int
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
			trendMD += "- [" + trend.Name + "](" + trend.Href + ") : " + trend.Language + "\n"
			if strings.TrimSpace(trend.Description) != "" {
				trendMD += "  - " + trend.Description + "\n"
			}
		} else {
			trendMD += "- [" + trend.Name + "](" + trend.Href + ")\n  - " + trend.Description + "\n"
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
		log.Fatalf("something wrong with create md file %s with %s", filename, err)
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
