package translator

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Translator interface {
	Translate(word string, fromLang Language, toLang Language) (Translation, error)
}

type reversoTranslator struct {
	contextEndpoint string
	headers         map[string]string
	maxMeanings     int
}

func NewReverso(maxMeanings int) *reversoTranslator {
	headers := make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
	headers["Origin"] = "https://context.reverso.net"
	headers["Accept-Language"] = "en-US,en;q=0.5"
	headers["X-Requested-With"] = "XMLHttpRequest"

	return &reversoTranslator{
		contextEndpoint: "https://context.reverso.net/translation/%s-%s/%s",
		maxMeanings:     maxMeanings,
		headers:         headers,
	}
}

func (rt reversoTranslator) Translate(word string, fromLang Language, toLang Language) (Translation, error) {
	endpoint := fmt.Sprintf(rt.contextEndpoint, fromLang, toLang, url.QueryEscape(word))

	req, _ := http.NewRequest("GET", endpoint, nil)
	for header, value := range rt.headers {
		req.Header.Add(header, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return Translation{}, err
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	return NewTranslation(word, fromLang, toLang, rt.parseMeanings(string(content)), rt.parseExamples(string(content))), nil
}

func (rt reversoTranslator) parseMeanings(apiResponse string) []Meaning {
	sectionPattern := regexp.MustCompile("id=\"top-results\"(.|\\s)+?<.+id=\"examples-content\"")
	res := sectionPattern.FindString(apiResponse)

	wordPattern := regexp.MustCompile("data-freq=\"(\\d+)\"(.|\\s)+?<span class=\"display-term\">([a-zA-Z ]+?)</span>( <span class=\"gender\">(\\w)</span>)?")

	matchedWords := wordPattern.FindAllStringSubmatch(res, -1)

	// sort words by descending usage frequency
	sort.Slice(matchedWords, func(i int, j int) bool {
		firstWordUsageFrequency, _ := strconv.Atoi(matchedWords[i][1])
		secondWordUsageFrequency, _ := strconv.Atoi(matchedWords[j][1])

		return firstWordUsageFrequency > secondWordUsageFrequency
	})

	var words []Meaning
	for _, matchedWord := range matchedWords {
		if rt.maxMeanings == len(words) {
			break
		}

		words = append(words, NewMeaning(matchedWord[3], matchedWord[5]))
	}

	return words
}

func (rt reversoTranslator) parseExamples(apiResponse string) []Example {
	sectionPattern := regexp.MustCompile("id=\"examples-content\"((.|\\s)+?)<section id=")
	res := sectionPattern.FindString(apiResponse)

	examplePattern := regexp.MustCompile("<div class=\"src ltr\">((.|\\s)+?)<div class=\"trg ltr\">((.|\\s)+?)<div class=\"options\"")
	// todo config
	examplesRaw := examplePattern.FindAllStringSubmatch(res, 5)

	var examples []Example
	for _, match := range examplesRaw {
		examples = append(examples, NewExample(sanitizeString(match[3]), sanitizeString(match[1])))
	}

	return examples
}

func sanitizeString(content string) string {
	pattern := regexp.MustCompile("(<.+?>)|(</.+?>)|\n")

	return strings.TrimSpace(pattern.ReplaceAllString(content, ""))
}
