from datetime import datetime
from pathlib import Path

import pandas as pd

from src.core.config import PROCESSED_DIR, RAW_DIR, SUBMISSIONS_DIR, ensure_data_dirs


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


def run_submit(
    predictions: str | Path = PROCESSED_DIR / "predictions.csv",
    sample_submission: str | Path = RAW_DIR / "sample_submission.csv",
    output: str | Path | None = None,
) -> Path:
    """Validate and save a submission CSV. Returns the output path."""
    ensure_data_dirs()

    pred_path = Path(predictions)
    sample_path = Path(sample_submission)

    if not pred_path.exists():
        raise FileNotFoundError(f"Missing predictions file: {pred_path}")
    if not sample_path.exists():
        raise FileNotFoundError(f"Missing sample submission: {sample_path}")

    pred_df = pd.read_csv(pred_path)
    sample_df = pd.read_csv(sample_path)

    validate_schema(pred_df, sample_df)

    if output is None:
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        output_path = SUBMISSIONS_DIR / f"submission_{timestamp}.csv"
    else:
        output_path = Path(output)

    output_path.parent.mkdir(parents=True, exist_ok=True)
    pred_df.to_csv(output_path, index=False)
    return output_path
