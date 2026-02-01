"""Data loading and preprocessing stubs."""

from __future__ import annotations

from pathlib import Path
import pandas as pd

from .config import RAW_DIR, PROCESSED_DIR


def load_raw(name: str) -> pd.DataFrame:
    """Load a raw CSV from data/raw/ by filename."""
    path = RAW_DIR / name
    if not path.exists():
        raise FileNotFoundError(f"Missing raw file: {path}")
    return pd.read_csv(path)


def save_processed(df: pd.DataFrame, name: str) -> Path:
    """Save a processed CSV to data/processed/."""
    path = PROCESSED_DIR / name
    path.parent.mkdir(parents=True, exist_ok=True)
    df.to_csv(path, index=False)
    return path
