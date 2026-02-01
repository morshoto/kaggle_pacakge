"""Central configuration for paths and common settings."""

from __future__ import annotations

import os
from pathlib import Path

ROOT_DIR = Path(__file__).resolve().parents[1]
DATA_DIR = Path(os.getenv("DATA_DIR", ROOT_DIR / "data"))
RAW_DIR = DATA_DIR / "raw"
INTERIM_DIR = DATA_DIR / "interim"
PROCESSED_DIR = DATA_DIR / "processed"
SUBMISSIONS_DIR = DATA_DIR / "submissions"
EXTERNAL_DIR = DATA_DIR / "external"

COMPETITION = os.getenv("COMPETITION", "")
SEED = int(os.getenv("SEED", "42"))


def ensure_data_dirs() -> None:
    """Create standard data directories if they do not exist."""
    for path in [RAW_DIR, INTERIM_DIR, PROCESSED_DIR, SUBMISSIONS_DIR, EXTERNAL_DIR]:
        path.mkdir(parents=True, exist_ok=True)
