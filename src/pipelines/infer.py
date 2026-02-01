from pathlib import Path

import pandas as pd

from src.core.config import PROCESSED_DIR, RAW_DIR, ensure_data_dirs


def run_infer(
    sample_submission: str | Path = RAW_DIR / "sample_submission.csv",
    output: str | Path = PROCESSED_DIR / "predictions.csv",
) -> Path:
    """Run inference and return predictions path."""
    ensure_data_dirs()

    sample_path = Path(sample_submission)
    if not sample_path.exists():
        raise FileNotFoundError(f"Missing sample submission: {sample_path}")

    df = pd.read_csv(sample_path)

    output_path = Path(output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    df.to_csv(output_path, index=False)
    return output_path
