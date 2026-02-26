import os
from pathlib import Path

ROOT_DIR = Path(__file__).resolve().parents[2]
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


def is_kaggle() -> bool:
    return Path("/kaggle/input").exists() or os.getenv("KAGGLE_URL_BASE") is not None


def resolve_raw_dir() -> Path:
    if os.getenv("RAW_DIR"):
        return Path(os.getenv("RAW_DIR")).expanduser()

    if is_kaggle() and COMPETITION:
        return Path("/kaggle/input") / COMPETITION

    data_dir = Path(os.getenv("DATA_DIR", ROOT_DIR / "data"))
    return data_dir / "raw"


RAW_DIR = resolve_raw_dir()
