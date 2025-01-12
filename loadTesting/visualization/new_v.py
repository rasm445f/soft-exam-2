import json
from datetime import datetime
import matplotlib.pyplot as plt
import pandas as pd
import gc
from collections import defaultdict

def process_json_in_chunks(file_path, chunk_size=1000):
    metrics = defaultdict(list)
    count = 0
    
    with open(file_path, "r") as file:
        chunk = []
        for line in file:
            try:
                entry = json.loads(line)
                if entry.get("type") == "Point":
                    metric = entry["metric"]
                    if metric in ["http_reqs", "http_req_duration", "http_req_blocked"]:
                        timestamp = datetime.fromisoformat(entry["data"]["time"].split("+")[0])
                        metrics[metric].append({
                            "time": timestamp,
                            "value": entry["data"]["value"]
                        })
                        count += 1
                        
                        # Process in chunks to free memory
                        if count % chunk_size == 0:
                            gc.collect()
                            
            except json.JSONDecodeError:
                continue
                
    return metrics

# Load and process data in chunks
file_path = "results.json"
metrics = process_json_in_chunks(file_path)

# Plot with memory optimization
plt.figure(figsize=(12, 8))

for metric, values in metrics.items():
    if values:
        # Convert to DataFrame one metric at a time
        df = pd.DataFrame(values)
        plt.plot(df["time"], df["value"], label=metric)
        
        # Clear memory
        del df
        gc.collect()

plt.title("K6 Load Testing Metrics")
plt.xlabel("Time")
plt.ylabel("Value")
plt.legend()
plt.grid(True)

# Save plot instead of showing it to free memory
plt.savefig('load_test_metrics.png')
plt.close()

# Clear remaining memory
gc.collect()