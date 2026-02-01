import argparse
import os
import sys
from pathlib import Path


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Download Kaggle competition data")
    parser.add_argument(
        "--competition",
        default=os.getenv("COMPETITION"),
        help="Kaggle competition slug (env: COMPETITION)",
    )
    parser.add_argument(
        "--dest",
        default=os.getenv("KAGGLE_DATA_DIR", "data/raw"),
        help="Destination directory (default: data/raw)",
    )
    parser.add_argument(
        "--no-unzip",
        action="store_true",
        help="Do not unzip downloaded files",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    if not args.competition:
        print(
            "ERROR: COMPETITION is not set. Use --competition or set COMPETITION env."
        )
        return 2

    try:
        from kaggle.api.kaggle_api_extended import KaggleApi
    except Exception as exc:  # pragma: no cover - best-effort import
        print("ERROR: Kaggle API is not installed. Run: pip install kaggle")
        print(f"Details: {exc}")
        return 3

    dest = Path(args.dest)
    dest.mkdir(parents=True, exist_ok=True)

    api = KaggleApi()
    api.authenticate()

    api.competition_download_files(
        args.competition,
        path=str(dest),
        unzip=not args.no_unzip,
    )

    print(f"Downloaded competition data to {dest}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
