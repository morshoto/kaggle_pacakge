"""Inference entry point stub."""

from __future__ import annotations

import argparse
from pathlib import Path

import pandas as pd

from .config import PROCESSED_DIR, RAW_DIR, ensure_data_dirs


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run inference (stub)")
    parser.add_argument(
        "--sample-submission",
        default=str(RAW_DIR / "sample_submission.csv"),
        help="Path to sample submission CSV",
    )
    parser.add_argument(
        "--output",
        default=str(PROCESSED_DIR / "predictions.csv"),
        help="Path to write predictions",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    ensure_data_dirs()

    sample_path = Path(args.sample_submission)
    if not sample_path.exists():
        raise FileNotFoundError(f"Missing sample submission: {sample_path}")

    df = pd.read_csv(sample_path)

    # Placeholder prediction: keep sample submission values as-is.
    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    df.to_csv(output, index=False)
    print(f"Wrote predictions to {output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
