import json
from datetime import datetime
import matplotlib.pyplot as plt
import pandas as pd

# Load JSON file
file_path = "results.json"  # Path to your JSON file
with open(file_path, "r") as file:
    data = [json.loads(line) for line in file]

# Parse metrics from JSON
metrics = {"http_reqs": [], "http_req_duration": [], "http_req_blocked": []}

for entry in data:
    if entry.get("type") == "Point":
        metric = entry["metric"]
        timestamp = datetime.fromisoformat(entry["data"]["time"].split("+")[0])
        value = entry["data"]["value"]

        if metric in metrics:
            metrics[metric].append({"time": timestamp, "value": value})

# Convert metrics to pandas DataFrames
dfs = {metric: pd.DataFrame(values) for metric, values in metrics.items()}

# Plot the metrics
plt.figure(figsize=(12, 8))
for metric, df in dfs.items():
    if not df.empty:
        plt.plot(df["time"], df["value"], label=metric)

plt.title("K6 Load Testing Metrics")
plt.xlabel("Time")
plt.ylabel("Value")
plt.legend()
plt.grid()
plt.tight_layout()

# Save the plot
plt.savefig("k6_metrics_visualization.png")
plt.show()
