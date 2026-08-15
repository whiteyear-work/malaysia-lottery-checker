// Package scraper retrieves lottery results from Toto Live.
package scraper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const totoLiveURL = "https://totolive.sportstoto.com.my/stoto/"

// SupremeTotoResult contains the latest Supreme Toto 6/58 result.
type SupremeTotoResult struct {
	DrawNumber string   `json:"draw_number"`
	DrawDate   string   `json:"draw_date"`
	Numbers    []string `json:"numbers"`
	Jackpot    string   `json:"jackpot"`
}

// FetchSupremeTotoResult fetches and validates the latest Supreme Toto 6/58 result.
func FetchSupremeTotoResult() (*SupremeTotoResult, error) {
	log.Printf("Fetching Toto Live result")
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodGet, totoLiveURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create Toto Live request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/128 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to fetch Toto Live result: %v", err)
		return nil, fmt.Errorf("fetch Toto Live result: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Toto Live returned HTTP status %s", resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse Toto Live HTML: %w", err)
	}
	result, err := parseSupremeToto(doc)
	if err != nil {
		log.Printf("Failed to parse Supreme Toto section: %v", err)
		return nil, err
	}
	log.Printf("Supreme Toto result fetched successfully")
	return result, nil
}

var numberPattern = regexp.MustCompile(`^\d+$`)

func parseSupremeToto(doc *goquery.Document) (*SupremeTotoResult, error) {
	var result *SupremeTotoResult
	doc.Find("body *").EachWithBreak(func(_ int, heading *goquery.Selection) bool {
		if normalize(heading.Text()) != "SUPREME TOTO 6/58" {
			return true
		}
		for ancestor := heading; ancestor.Length() > 0 && result == nil; ancestor = ancestor.Parent() {
			candidate := extractCandidate(ancestor)
			if len(candidate.Numbers) == 6 {
				result = candidate
			}
		}
		return result == nil
	})
	if result == nil {
		return nil, fmt.Errorf("unable to extract Supreme Toto 6/58 result")
	}
	return result, nil
}

func extractCandidate(section *goquery.Selection) *SupremeTotoResult {
	result := &SupremeTotoResult{}
	section.Find("*").Each(func(_ int, item *goquery.Selection) {
		if item.Children().Length() != 0 {
			return
		}
		text := strings.TrimSpace(item.Text())
		upper := strings.ToUpper(text)
		switch {
		case strings.Contains(upper, "JACKPOT"):
			value := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), ":"))
			if value != ":" {
				result.Jackpot = value
			}
		case strings.Contains(upper, "RM"):
			result.Jackpot = text
		case strings.Contains(upper, "DRAW") && result.DrawNumber == "":
			result.DrawNumber = text
		case strings.Contains(upper, "DATE") && result.DrawDate == "":
			result.DrawDate = text
		case numberPattern.MatchString(text) && len(result.Numbers) < 6:
			result.Numbers = append(result.Numbers, text)
		}
	})
	return result
}

func normalize(value string) string {
	return strings.Join(strings.Fields(strings.ToUpper(value)), " ")
}
