# AGENTS.md

## Project

Malaysia Lottery Checker is a small personal Go application.

The application will:

1. Read purchased Toto 6D numbers from XML.
2. Scrape the latest Toto 6D result.
3. Compare the winning digits with each purchased number.
4. Count matching digits regardless of position.
5. Write the latest result back to XML.
6. Display results through a basic HTML page.
7. Send a Telegram notification after scheduled checks.

## Technology

- Go
- Docker
- Docker Compose
- Standard Go HTML templates
- XML storage
- Basic CSS
- Telegram Bot API
- Ubuntu cron in production

## Development environment

- Windows
- Docker Desktop
- VS Code

Do not require Go to be installed directly on Windows.

## Production environment

- Ubuntu
- Docker Compose
- Linux cron

## Coding guidelines

- Keep the project simple.
- Prefer Go standard-library packages where practical.
- Do not add a database.
- Do not add Node.js, React, Redis, queues, or microservices.
- Keep lottery numbers as strings to preserve leading zeroes.
- Use relative file paths such as `data/numbers.xml`.
- Never hardcode Telegram secrets.
- Return and handle errors rather than ignoring them.
- Format Go code using `gofmt`.