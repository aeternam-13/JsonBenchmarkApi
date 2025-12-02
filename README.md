## 📊 Visualization Setup (Python)

To generate the performance charts (`benchmark_analysis.png`) from the CSV data, a python program was created, it requires a venv

### Prerequisites
- Python 3.x
- `pip`

### Create the Virtual Environment
We use a virtual environment named `stadistics` to manage dependencies.

```bash
python -m venv stadistics
source chart/bin/activate
```

To activate the enviroment and get the dependencies

```bash
source chart/bin/activate
pip install -r requirements.txt
```
Finally run the program

```bash
python stadistics.py
```