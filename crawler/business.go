package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func buildGoogleUrls(searchTerm, countryCode, languageCode string, pages, count int) ([]string, error) {
	toScrape := []string{}
	searchTerm = strings.Trim(searchTerm, "")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)

	if googleBase, found := googleDomains[countryCode]; found {
		for i := 0; i < pages; i++ {
			start := i * count
			scrapeURL := fmt.Sprintf("%s%s&num=%d&hl=%s&start=%d&filter=0 ", googleBase, searchTerm, count, languageCode, start)

			toScrape = append(toScrape, scrapeURL)
		}
	} else {
		err := fmt.Errorf("country (%s) is currently not supported", countryCode)
		return nil, err
	}
	return toScrape, nil
}

func GoogleScrape(searchTerm, countryCode, languageCode string, proxyString interface{}, pages, count int) ([]SearchResult, error) {
	results := []SearchResult{}
	resultCounter := 0

	googlePages, err := buildGoogleUrls(searchTerm, countryCode, languageCode, pages, count)
	if err != nil {
		return nil, err
	}

	for _, page := range googlePages {
		res, err := scrapeClientRequest(page, proxyString)
		if err != nil {
			return nil, err
		}

		data, err := googleResultParsing(res, resultCounter)
		if err != nil {
			return nil, err
		}

		resultCounter += len(data)
		results = append(results, data...)
	}

	return results, nil
}

func scrapeClientRequest(searchURL string, proxyString interface{}) (*http.Response, error) {
	baseClient := getScrapeClient(proxyString)
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set("User-Agent", randomUserAgent())

	res, err := baseClient.Do(req)
	if res.StatusCode != 200 {
		err := fmt.Errorf("scraper received a non-200 status code suggesting a ban")
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func googleResultParsing(response *http.Response, rank int) ([]SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}
	results := []SearchResult{}
	sel := doc.Find("div.g")
	rank++
	for i := range sel.Nodes {
		item := sel.Eq(i)
		linkTag := item.Find("a")
		link, _ := linkTag.Attr("href")

		link = strings.Trim(link, "")

		if link != "" && link != "#" && !strings.HasPrefix(link, "/") {
			result := SearchResult{
				ResultRank:  rank,
				ResultURL:   link,
				ResultTitle: item.Find("h3").Text(),
				ResultDesc:  item.Find("-webkit-line-clamp:2").Text(),
			}

			results = append(results, result)
			rank++
		}
	}

	return results, err
}

func getScrapeClient(proxyString interface{}) *http.Client {
	switch v := proxyString.(type) {

	case string:
		proxyURL, _ := url.Parse(v)
		return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	default:
		return &http.Client{}
	}

}
