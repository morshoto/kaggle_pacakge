import numpy as np


def fit(X: np.ndarray, y: np.ndarray) -> dict:
    """Return a placeholder model artifact."""
    return {"mean": float(np.mean(y)) if y.size else 0.0}


def predict(model: dict, X: np.ndarray) -> np.ndarray:
    """Return constant predictions based on the fitted mean."""
    mean = float(model.get("mean", 0.0))
    return np.full((X.shape[0],), mean, dtype=float)
