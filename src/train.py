"""Training entry point stub."""

from __future__ import annotations

import argparse
from pathlib import Path

from .config import PROCESSED_DIR, ensure_data_dirs


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
    ensure_data_dirs()

    output = Path(args.output)
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text("placeholder model artifact\n")
    print(f"Wrote model artifact to {output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
