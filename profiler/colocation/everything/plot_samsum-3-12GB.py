import pandas as pd
import matplotlib.pyplot as plt
import ast

data = pd.read_csv("samsum-3-12GB.csv")
data['memory'] = data['memory1'] + data['memory2'] + data['memory3']
data['compute'] = data['compute1'] + data['compute2'] + data['compute3']

data = data[data['load'] == 128]
data = data[["compute", "utilization", "latencies"]]

compute_values = data['compute'].unique()

colocated3 = pd.read_csv("../samsum-3-colocation/results.csv")
colocated3["memory"] = colocated3["memory1"] + colocated3["memory2"] + colocated3["memory3"]
colocated3["compute"] = colocated3["compute1"] + colocated3["compute2"] + colocated3["compute3"]

colocated3o = pd.read_csv("../samsum-3-overprovisioned/results.csv")
colocated3o["memory"] = colocated3o["memory1"] + colocated3o["memory2"] + colocated3o["memory3"]
colocated3o["compute"] = colocated3o["compute1"] + colocated3o["compute2"] + colocated3o["compute3"]

colocated3 = colocated3[["memory", "compute", "load", "latency", "throughput"]]
colocated3o = colocated3o[["memory", "compute", "load", "latency", "throughput"]]
dataset3 = pd.concat([colocated3, colocated3o], ignore_index=True)

plt.figure(figsize=(30, 10))

plt.subplot(1, 3, 1)
for com in compute_values:
    subset = dataset3[(dataset3["compute"] == com) & (dataset3["memory"] == 12)]
    if com == 300:
        plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
    elif com > 100:
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores (overprovisioned): {com/3}%, {com/3}% and {com/3}%')
    else:
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores: {com/3}%, {com/3}% and {com/3}%')

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
plt.xlabel('Total GPU cores (Each workload has 33% total cores allocated)')
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

plt.suptitle("Throughput, Utilization and Latency distribution for 3 colocated workloads running with 6GB memory each")
plt.savefig("samsum-3-12GB.png")