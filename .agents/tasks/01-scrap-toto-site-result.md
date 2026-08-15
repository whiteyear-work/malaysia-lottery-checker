# Agent Task: Scrape Supreme Toto 6/58 Result

## Objective

Implement the first lottery scraper for the Malaysia Lottery Checker.

The scraper must fetch the latest **Supreme Toto 6/58** result from:

`https://totolive.sportstoto.com.my/stoto/`

For this task, only support **Supreme Toto 6/58**.

Do not implement Telegram, XML persistence, matching logic, cron scheduling, or other lottery types yet.

---

## Expected Result

From the Supreme Toto 6/58 section, extract:

* Draw number if available
* Draw date if available
* Six winning numbers
* Jackpot amount if available

Example based on the current page layout:

```text
SUPREME TOTO 6/58

25
9
26
15
22
19

Jackpot:
RM 21,317,577.20
```

The six winning numbers must be returned in the same order shown on the website.

---

## Project Structure

Create a scraper package:

```text
scraper/
└── toto.go
```

Keep `main.go` responsible only for starting the application and testing/invoking the scraper.

Do not put all scraping logic directly inside `main.go`. May refer `exmaple\sample.txt` for exmaple of html structure

---

## Result Model

Create an appropriate Go struct for the scraped result.

Example:

```go
type SupremeTotoResult struct {
    DrawNumber string
    DrawDate   string
    Numbers    []string
    Jackpot    string
}
```

Numbers must be stored as strings.

Do not use integers because lottery values may eventually require preserving formatting.

---

## Scraper Function

Create a function similar to:

```go
func FetchSupremeTotoResult() (*SupremeTotoResult, error)
```

The function should:

1. Send an HTTP GET request to the Toto Live website.
2. Use a reasonable HTTP timeout.
3. Check that the HTTP response status is successful.
4. Parse the HTML response.
5. Find the `SUPREME TOTO 6/58` section.
6. Extract exactly six winning numbers.
7. Extract jackpot amount if available.
8. Extract draw number and draw date if available.
9. Validate the result.
10. Return an error if valid Supreme Toto data cannot be extracted.

---

## HTTP Request

Use Go's HTTP client.

Configure a timeout such as:

```go
15 * time.Second
```

The website may reject requests that do not resemble browser traffic.

Use normal browser-like request headers where needed, for example:

```text
User-Agent
Accept
Accept-Language
```

Do not add unnecessary browser automation unless normal HTTP parsing cannot retrieve the result.

---

## HTML Parsing

Use:

```text
github.com/PuerkitoBio/goquery
```

Install it through Go modules.

Do not parse HTML using regular expressions.

The scraper must inspect the actual HTML structure and determine suitable selectors.

Do not guess selectors from the screenshot alone.

Selectors should be kept clear and easy to update because the Toto page HTML may change in the future.

---

## Finding the Supreme Toto Section

Do not rely only on absolute element positions such as:

```text
table:nth-child(5)
```

Prefer locating the relevant section using its visible heading:

```text
SUPREME TOTO 6/58
```

Then traverse nearby DOM elements to find:

* six numbers
* jackpot

The scraper should avoid accidentally returning numbers from:

```text
Power Toto 6/55
Star Toto
Toto 4D
```

Only data belonging to **Supreme Toto 6/58** is valid for this task.

---

## Validation

Before returning the result, validate:

```text
len(Numbers) == 6
```

Each winning number should contain only numeric characters.

The normal Supreme Toto number range is expected to be compatible with 6/58, but avoid making unnecessary assumptions unless required.

If six valid numbers cannot be found, return an error instead of silently returning incomplete data.

Example:

```go
return nil, fmt.Errorf("unable to extract Supreme Toto 6/58 result")
```

---

## Logging

Use Go's existing logging approach.

Log important scraper events such as:

```text
Fetching Toto Live result
Supreme Toto result fetched successfully
Failed to fetch Toto Live result
Failed to parse Supreme Toto section
```

Do not log the entire HTML response.

---

## Error Handling

Handle at minimum:

* HTTP timeout
* DNS/network failure
* non-success HTTP response
* unexpected HTML structure
* Supreme Toto section missing
* fewer or more than six result numbers
* invalid number values

Return errors to the caller.

Do not panic inside scraper functions.

---

## Temporary Testing Through `main.go`

For now, allow the result to be printed when the application starts or through a simple temporary test route.

Preferred temporary route:

```text
GET /test/supreme-toto
```

Example JSON response:

```json
{
    "draw_number": "",
    "draw_date": "",
    "numbers": [
        "25",
        "9",
        "26",
        "15",
        "22",
        "19"
    ],
    "jackpot": "RM 21,317,577.20"
}
```

This route is only for development/testing.

Do not build the final frontend yet.

---

## Docker Compatibility

The scraper must work inside the existing Docker development environment.

Do not require Go to be installed directly on Windows.

After adding dependencies, ensure these work:

```bash
go mod tidy
go run .
```

The application must continue working with Air auto reload inside Docker.

---

## Coding Guidelines

Follow the existing project guidelines:

* Keep code simple.
* Prefer standard Go patterns.
* Keep `main.go` thin.
* Put scraping logic inside the scraper package.
* Add a short comment for each exported function.
* Handle external HTTP calls with proper error handling and logging.
* Do not ignore returned errors.
* Format code using `gofmt`.
* Do not add database-related code.
* Do not add Redis.
* Do not add Node.js.
* Do not implement Telegram yet.
* Do not implement scheduled jobs yet.

---

## Deliverables

Implement:

```text
scraper/toto.go
```

Update:

```text
go.mod
go.sum
main.go
```

if necessary.

The completed task should allow:

```text
GET /test/supreme-toto
```

to return the latest Supreme Toto 6/58 result retrieved from the Toto Live website.

---

## Acceptance Criteria

The task is complete when:

1. The application successfully requests the Toto Live page.
2. The scraper identifies the correct Supreme Toto 6/58 section.
3. Exactly six winning numbers are returned.
4. Numbers remain in website display order.
5. Jackpot is returned when present.
6. Missing or changed website structure produces a clear error.
7. The scraper does not accidentally return Power Toto or another game result.
8. The test endpoint works from the Docker development environment.
9. `go mod tidy` completes successfully.
10. Existing application startup and home page continue to work.
11. Show implementation plan and wait for approval before starting task. 
