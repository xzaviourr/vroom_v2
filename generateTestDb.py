import requests
import json
import pandas as pd

profile = pd.read_csv("profiler/colocation/results/cnn-1-fixed-memory-all.csv")

url = "http://localhost:8083/insert"
headers = {"Content-Type": "application/json"}

for index, row in profile[profile["arrival_rate"] == 12].iterrows():
    print(row["memory1"])

    minimum_latency = profile[(profile['memory1'] == row["memory1"]) & 
            (profile['compute1'] == row["compute1"]) & 
             (profile['arrival_rate'] == 2)]["average_latency"].iloc[0]
    
    average_latency = profile[(profile['memory1'] == row["memory1"]) & 
            (profile['compute1'] == row["compute1"]) & 
             (profile['arrival_rate'] == 8)]["average_latency"].iloc[0]
    
    max_latency = profile[(profile['memory1'] == row["memory1"]) & 
            (profile['compute1'] == row["compute1"]) & 
             (profile['arrival_rate'] == 16)]["average_latency"].iloc[0]
    
    data = {
        "task-identifier": "text-summarization",
        "gpu-memory": row["memory1"],
        "gpu-cores": row["compute1"],
        "image": "synergcseiitb/bart-large-cnn-text_summarization",
        "startup-latency": row["startup_time"],
        "min-latency": minimum_latency,
        "mean-latency": average_latency,
        "max-latency": max_latency,
        "accuracy": 85,
        "batch-size": 64,
        "end-point": "/summarize",
        "port": 4444,
        "capacity": row["throughput"]
    }

    print(data)

    response = requests.post(url, headers=headers, data=json.dumps(data))

    print(f"Status Code: {response.status_code}")
    print(f"Response: {response.json()}")
