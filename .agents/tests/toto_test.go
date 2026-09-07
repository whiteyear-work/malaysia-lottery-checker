package tests

import (
	"malaysia-lottery-checker/scraper"
	"os"
	"strings"
	"testing"
)

func TestSamples(t *testing.T) {
	for _, tc := range []struct{ file, numbers, jackpot, draw string }{
		{"sample1.txt", "25,9,26,15,22,19", "RM 21,317,577.20", "6171/26"},
		{"sample2.txt", "24,2,23,35,34,41", "RM 10,241,508.57", "6182/26"},
	} {
		t.Run(tc.file, func(t *testing.T) {
			data, err := os.ReadFile("../../example/" + tc.file)
			if err != nil {
				t.Fatal(err)
			}
			got, err := scraper.ParseSupremeTotoResult(strings.NewReader(string(data)))
			if err != nil {
				t.Fatal(err)
			}
			if strings.Join(got.Numbers, ",") != tc.numbers || got.Jackpot != tc.jackpot || got.DrawNumber != tc.draw || got.DrawDate == "" {
				t.Fatalf("unexpected result: %+v", got)
			}
		})
	}
}

func TestInvalidNumbers(t *testing.T) {
	for _, values := range []string{"1,2,3,4,5", "1,2,3,4,5,6,7", "1,2,3,4,5,**"} {
		html := `<table><tr><td>SUPREME TOTO 6/58</td></tr></table><table class="supremetoto"><tr><td>` + strings.ReplaceAll(values, ",", "</td><td>") + `</td></tr></table><table><tr><td>POWER TOTO 6/55</td></tr><tr><td>1</td><td>2</td><td>3</td><td>4</td><td>5</td><td>6</td></tr></table>`
		if got, err := scraper.ParseSupremeTotoResult(strings.NewReader(html)); err == nil {
			t.Fatalf("accepted invalid result: %+v", got)
		}
	}
}
