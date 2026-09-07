// Package scraper retrieves lottery results from Toto Live.
package scraper

import (
	"fmt"
	"io"
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

	result, err := ParseSupremeTotoResult(resp.Body)
	if err != nil {
		log.Printf("Failed to parse Supreme Toto section: %v", err)
		return nil, err
	}
	log.Printf("Supreme Toto result fetched successfully")
	return result, nil
}

var numberPattern = regexp.MustCompile(`^\d+$`)

// ParseSupremeTotoResult extracts a Supreme result from Toto Live HTML.
func ParseSupremeTotoResult(reader io.Reader) (*SupremeTotoResult, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("parse Toto Live HTML: %w", err)
	}
	// The heading and result are separate sibling tables in Toto Live HTML.
	headings := doc.Find("td").FilterFunction(func(_ int, cell *goquery.Selection) bool {
		return normalize(cell.Text()) == "SUPREME TOTO 6/58"
	})
	if headings.Length() != 1 {
		return nil, fmt.Errorf("expected one Supreme Toto 6/58 heading, found %d", headings.Length())
	}
	table := headings.Closest("table").Next()
	if !table.Is("table.supremetoto") {
		return nil, fmt.Errorf("Supreme Toto 6/58 result table missing after heading")
	}
	rows := table.ChildrenFiltered("tbody").ChildrenFiltered("tr").AddSelection(table.ChildrenFiltered("tr"))
	cells := rows.First().ChildrenFiltered("td")
	if cells.Length() != 6 {
		return nil, fmt.Errorf("Supreme Toto 6/58 requires six numbers, found %d", cells.Length())
	}
	result := &SupremeTotoResult{}
	for i := 0; i < cells.Length(); i++ {
		number := strings.TrimSpace(cells.Eq(i).Text())
		if !numberPattern.MatchString(number) {
			return nil, fmt.Errorf("invalid Supreme Toto number at position %d: %q", i+1, number)
		}
		result.Numbers = append(result.Numbers, number)
	}
	rows.Slice(1, rows.Length()).Each(func(_ int, row *goquery.Selection) {
		fields := row.ChildrenFiltered("td")
		if normalize(fields.First().Text()) == "JACKPOT" {
			result.Jackpot = strings.Join(strings.Fields(fields.Eq(1).Text()), " ")
		}
	})
	table.Closest("[id^='resultTable']").Find("span").Each(func(_ int, span *goquery.Selection) {
		label, value, ok := strings.Cut(span.Text(), ":")
		if !ok {
			return
		}
		switch normalize(label) {
		case "DRAW NO":
			result.DrawNumber = strings.TrimSpace(value)
		case "DRAW DATE":
			result.DrawDate = strings.TrimSpace(value)
		}
	})
	return result, nil
}

func normalize(value string) string {
	return strings.Join(strings.Fields(strings.ToUpper(value)), " ")
}
