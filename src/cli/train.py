import argparse

from src.core.config import PROCESSED_DIR
from src.pipelines.train import run_train


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Train model (stub)")
    parser.add_argument(
        "--output",
        default=str(PROCESSED_DIR / "model.txt"),
        help="Path to write a placeholder model artifact",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    output = run_train(output=args.output)
    print(f"Wrote model artifact to {output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
