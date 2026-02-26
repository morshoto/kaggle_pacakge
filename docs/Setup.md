# Project Setup

This document describes how to set up the environment and download Kaggle competition data for this repository.

## Prerequisites

- Python 3.10+
- `uv` installed
- A Kaggle account
- You have **joined the competition** and **accepted the rules** in the Kaggle UI

## 1) Configure Kaggle credentials

This project uses **kagglehub** for downloading competition data (recommended), with a fallback to the legacy Kaggle CLI.

### Create a Kaggle API token

1. Go to your Kaggle Account settings page.
2. Create an **API Token** (recommended tokens are supported by kagglehub).
3. Copy the token value.

### Create `.env`

Create a `.env` file at the repository root (or copy from `.env.example` if present):

```bash
# kagglehub token (recommended)
KAGGLE_API_TOKEN="KGAT_..."

# legacy Kaggle CLI credentials (optional fallback)
KAGGLE_USERNAME="your_username"
KAGGLE_KEY="your_key"

COMPETITION="playground-series-s6e2"
```

> Notes:
>
> Notes:
>
> - `KAGGLE_API_TOKEN` is required for kagglehub.
> - `KAGGLE_USERNAME` and `KAGGLE_KEY` are only used if kagglehub fails and the legacy CLI is used.
> - The competition slug is the part after `/competitions/` in the Kaggle URL.

## 2) Install dependencies

```bash
uv sync
```

(If you add packages, use `uv add <pkg>` and commit the updated `pyproject.toml` and `uv.lock`.)

## 3) Load env vars and download the data

`uv` does not automatically load `.env`, so we export variables before running.

```bash
set -a
source .env
set +a

uv run python -m scripts.download --competition playground-series-s6e2 --dest data/raw
```

Expected output:

- Downloads an archive into the kagglehub cache directory
- Extracts files
- Writes data into `data/raw`

## 4) Verify downloaded files

```bash
ls -la data/raw
```

You should see files such as `train.csv`, `test.csv`, and `sample_submission.csv` (exact names depend on the competition).

## Troubleshooting

### 401 Unauthorized

If you see `401 Client Error: Unauthorized`, check the following:

1. You are logged into the correct Kaggle account (same account as your token).
2. You have joined the competition and accepted the rules in the browser.
3. `KAGGLE_API_TOKEN` is set in your current shell:

```bash
echo "$KAGGLE_API_TOKEN" | head -c 6; echo
```

### Legacy Kaggle CLI / kaggle.json

This repository prefers **kagglehub**. The legacy `kaggle` CLI / `kaggle.json` credentials may not work with newer token formats (e.g., `KGAT_...`).
