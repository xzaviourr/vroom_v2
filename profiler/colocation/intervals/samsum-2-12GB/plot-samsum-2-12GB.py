import pandas as pd
import matplotlib.pyplot as plt
import ast

data = pd.read_csv("samsum-2-12GB-prometheus.csv")
data['memory'] = data['memory1'] + data['memory2']
data['compute'] = data['compute1'] + data['compute2']
data = data[["memory", "compute", "load", "throughput", "utilization", "latencies"]]

compute_values = data['compute'].unique()


plt.figure(figsize=(30, 10))

# Throughput Plot
plt.subplot(1, 3, 1)
for com in compute_values:
    subset = data[(data["compute"] == com) & (data["memory"] == 12)]
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
plt.xticks([1,2,3,4,5,6,7,8,9,10])

# =================================================================================================================
# Latency Distribution Plot

plt.subplot(1, 3, 2)
dataset = []
for com in compute_values:
    latencies = data[(data['compute'] == com) & (data['memory'] == 12) & (data['load'] == 10)]["latencies"].iloc[0]
    latencies = ast.literal_eval(latencies)  
    dataset.append(latencies)

plt.boxplot(dataset)
plt.xticks(range(1, 11), compute_values)  
plt.xlabel('Total GPU cores percentage (Divided equally among the 2 workloads)')
plt.ylabel('Response Time (ms)')
plt.title(f'Distribution of Response Times (Load = 10 Req/s)')
plt.grid(True)

# =================================================================================================================
# GPU Utilization Plot

plt.subplot(1, 3, 3)
for com in compute_values:
    utilization = data[(data['compute'] == com) & (data['memory'] == 12) & (data['load'] == 10)]["utilization"].iloc[0]

    utilization = ast.literal_eval(utilization)    
    df = pd.DataFrame(utilization, columns=['Time', 'GPU_Utilization'])
    df['Time'] = df['Time'] - df['Time'].iloc[0]  # Normalize the time
    plt.plot(df['Time'], df['GPU_Utilization'], label=f"GPU cores: {com}%")

plt.xlabel('Time (s)')
plt.ylabel('GPU Utilization (%)')
plt.title(f'GPU Utilization Over Time (Load = 10 Req/s)')
plt.grid(True)
plt.legend()

plt.suptitle("Throughput, Utilization and Latency distribution for 2 colocated workloads running with 6GB memory each")
plt.savefig("samsum-2-12GB-prometheus.png")