"""Submission generation with schema checks."""

from __future__ import annotations

import argparse
from datetime import datetime
from pathlib import Path

import pandas as pd

from .config import PROCESSED_DIR, RAW_DIR, SUBMISSIONS_DIR, ensure_data_dirs


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate and create submission CSV")
    parser.add_argument(
        "--predictions",
        default=str(PROCESSED_DIR / "predictions.csv"),
        help="Predictions CSV to validate",
    )
    parser.add_argument(
        "--sample-submission",
        default=str(RAW_DIR / "sample_submission.csv"),
        help="Sample submission CSV for schema reference",
    )
    parser.add_argument(
        "--output",
        default="",
        help="Output CSV path (default: data/submissions/submission_YYYYmmdd_HHMMSS.csv)",
    )
    return parser.parse_args()


def validate_schema(pred_df: pd.DataFrame, sample_df: pd.DataFrame) -> None:
    if list(pred_df.columns) != list(sample_df.columns):
        raise ValueError(
            "Column mismatch. Expected columns in order: "
            f"{list(sample_df.columns)} but got {list(pred_df.columns)}"
        )

    if pred_df.isnull().any().any():
        missing = pred_df.isnull().sum()
        missing = missing[missing > 0]
        raise ValueError(f"Predictions contain missing values: {missing.to_dict()}")

    # dtype checks can be added per competition if needed.


def main() -> int:
    args = parse_args()
    ensure_data_dirs()

    pred_path = Path(args.predictions)
    sample_path = Path(args.sample_submission)

    if not pred_path.exists():
        raise FileNotFoundError(f"Missing predictions file: {pred_path}")
    if not sample_path.exists():
        raise FileNotFoundError(f"Missing sample submission: {sample_path}")

    pred_df = pd.read_csv(pred_path)
    sample_df = pd.read_csv(sample_path)

    validate_schema(pred_df, sample_df)

    if args.output:
        output = Path(args.output)
    else:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output = SUBMISSIONS_DIR / f"submission_{timestamp}.csv"

    output.parent.mkdir(parents=True, exist_ok=True)
    pred_df.to_csv(output, index=False)
    print(f"Saved submission to {output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
