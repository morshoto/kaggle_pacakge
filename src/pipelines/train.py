from pathlib import Path

from src.core.config import PROCESSED_DIR, ensure_data_dirs


def run_train(output: str | Path = PROCESSED_DIR / "model.txt") -> Path:
    """Run training and return the model artifact path."""
    ensure_data_dirs()

    output_path = Path(output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text("placeholder model artifact\n")
    return output_path
