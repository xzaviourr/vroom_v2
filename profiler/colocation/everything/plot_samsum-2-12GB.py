import pandas as pd
import matplotlib.pyplot as plt
import ast

data = pd.read_csv("results.csv")
data['memory'] = data['memory1'] + data['memory2']
data['compute'] = data['compute1'] + data['compute2']

data = data[data['load'] == 1024]
data = data[["compute", "utilization", "latencies"]]

compute_values = data['compute'].unique()

colocated2 = pd.read_csv("../samsum-2-colocation/results.csv")
colocated2["memory"] = colocated2["memory1"] + colocated2["memory2"]
colocated2["compute"] = colocated2["compute1"] + colocated2["compute2"]
colocated2o = pd.read_csv("../samsum-2-overprovisioned/results.csv")
colocated2o["memory"] = colocated2o["memory1"] + colocated2o["memory2"]
colocated2o["compute"] = colocated2o["compute1"] + colocated2o["compute2"]
colocated2 = colocated2[["memory", "compute", "load", "latency", "throughput"]]
colocated2o = colocated2o[["memory", "compute", "load", "latency", "throughput"]]
dataset2 = pd.concat([colocated2, colocated2o], ignore_index=True)

plt.figure(figsize=(30, 10))

plt.subplot(1, 3, 1)
for com in compute_values:
    subset = dataset2[(dataset2["compute"] == com) & (dataset2["memory"] == 12)]
    if com == 200:
        plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
    elif com > 100:
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores (overprovisioned): {com/2}% and {com/2}%')
    else:
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores: {com/2}% and {com/2}%')

plt.title(f"Distribution of Throughput")
plt.legend()
plt.grid(True)
plt.xlabel("Load (Requests queued at once)")
plt.ylabel("Throughput (Reqs / sec)")
plt.xticks([2, 4, 8, 16, 32, 64, 128])

# =================================================================================================================

plt.subplot(1, 3, 2)
dataset = []
for com in compute_values:
    latencies = data[data['compute'] == com]["latencies"].iloc[0]
    latencies = ast.literal_eval(latencies)  
    dataset.append(latencies)

plt.boxplot(dataset)
plt.xticks(range(1, 11), compute_values)  
plt.xlabel('Total GPU cores (Each workload has 50% total cores allocated)')
plt.ylabel('Response Time (ms)')
plt.title(f'Distribution of Response Times (Load = 128 Req/s)')
plt.grid(True)

# =================================================================================================================

plt.subplot(1, 3, 3)
for com in compute_values:
    utilization = data[data['compute'] == com]["utilization"].iloc[0]

    utilization = ast.literal_eval(utilization)    
    df = pd.DataFrame(utilization, columns=['Time', 'GPU_Utilization'])
    df['Time'] = df['Time'] - df['Time'].iloc[0]  # Normalize the time
    plt.plot(df['Time'], df['GPU_Utilization'], label=f"GPU cores: {com}%")

plt.xlabel('Time (s)')
plt.ylabel('GPU Utilization (%)')
plt.title(f'GPU Utilization Over Time (Load = 128 Req/s)')
plt.grid(True)
plt.legend()

plt.suptitle("Throughput, Utilization and Latency distribution for 2 colocated workloads running with 6GB memory each")
plt.savefig("test.png")