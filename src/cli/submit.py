import argparse

from src.core.config import PROCESSED_DIR, RAW_DIR
from src.pipelines.submit import run_submit


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


def main() -> int:
    args = parse_args()
    output = run_submit(
        predictions=args.predictions,
        sample_submission=args.sample_submission,
        output=args.output if args.output else None,
    )
    print(f"Saved submission to {output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
