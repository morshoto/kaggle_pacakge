# get_discussion

CLI for downloading Kaggle discussion threads to Markdown files.

## Requirements

- Go 1.22+

## Usage

Initialize `.env` file from `.env.example` file to read competition. From repo root:

```bash
# Single discussion
go run ./cli/get_discussion --link "https://www.kaggle.com/discussion/12345"
# Same command using the module import path
go run github.com/shotomorisaki/kaggle_pacakge/cli/get_discussion --link "https://www.kaggle.com/discussion/12345"
# Listing (default 10 discussions)
go run ./cli/get_discussion --sort hotness --time-filter last_7_days
# Listing with more than 10 discussions
go run ./cli/get_discussion --sort hotness --time-filter last_7_days --limit 25
# Listing with all discussions
go run ./cli/get_discussion --sort hotness --time-filter last_7_days --all
# Competition listing
go run ./cli/get_discussion --sort most_votes --time-filter last_30_days
```

## Flags

- `--link`: Download a single discussion by URL.
- `--sort`: `hotness`, `recent_comments`, `recently_posted`, `most_votes`, `most_comments`.
- `--time-filter`: `last_30_days`, `last_7_days`, `today`.
- `--output-dir`: Output directory for Markdown files (default `discussion`).
- `--limit`: Max discussions to download when listing (default `10`).
- `--all`: Download all discussions (ignores `--limit`).
- `--delay`: Delay in seconds between requests (default `0.5`).
- `--verbose`: Enable verbose logging.

## Environment

- `COMPETITION`: If set, fetches discussions from a specific Kaggle competition forum.

## Output

Markdown files are written to `--output-dir` with YAML front matter:

- `title`
- `link`
- `author`
- `comments`
- `published_date`
