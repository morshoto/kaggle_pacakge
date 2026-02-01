import argparse

from src.core.config import PROCESSED_DIR, RAW_DIR
from src.pipelines.infer import run_infer


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
    output = run_infer(
        sample_submission=args.sample_submission,
        output=args.output,
    )
    print(f"Wrote predictions to {output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
