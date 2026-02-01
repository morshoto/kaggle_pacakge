"""Feature engineering stubs."""

from __future__ import annotations

import pandas as pd


def build_features(df: pd.DataFrame) -> pd.DataFrame:
    """Placeholder feature builder. Extend per competition."""
    return df.copy()
