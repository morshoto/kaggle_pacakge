# kaggle_comp_temp

### Description

Develop machine learning models to predict. The goal is to improve understanding.

```bash
├── README.md
├── data           <---- Data directory
├── docs           <---- Documents, logs,
│   ├── Log.md     <---- Day tracking work
│   ├── Paper.md   <---- Paper research
│   └── Scoring.md <---- Score tracking table
├── paper          <---- Papers to read, get inspired
├── nb             <---- Created on jupyter notebook
├── nb_download    <---- Public notebook from kaggle
└── src            <---- Globally used functions
```

### Dataset

The dataset provided for this competition consists of.

| Name | Detail | Size | Link     |
| ---- | ------ | ---- | -------- |
| name |        |      | [Link]() |

### Quickstart

```bash
# 1) set env vars (or copy .env.example to .env and edit)
export COMPETITION=your-competition
export KAGGLE_USERNAME=your_name
export KAGGLE_KEY=your_key

# 2) download data to data/raw
make setup

# 3) train / predict / submit
make train
make predict
make submit
```

### Data layout (standard)

```bash
data/
  raw/           # Kaggle downloads
  interim/       # intermediate artifacts
  processed/     # training/inference ready data
  submissions/   # generated submission CSVs
  external/      # external datasets
```

### Notebook naming

Use numbered notebooks in `nb/` to keep ordering clear, e.g. `00_eda.ipynb`, `01_baseline.ipynb`.
