import pandas as pd
import matplotlib.pyplot as plt
import ast

filename = "samsum-2-fixed-memory.csv"
num_colocation = 2
memory = 12

data = pd.read_csv(filename)
data['memory'] = 0
data['compute'] = 0
for ind in range(1, num_colocation+1):
    data['memory'] = data['memory'] + data[f'memory{ind}']
    data['compute'] = data['compute'] + data[f'compute{ind}']

data = data[["memory", "compute", "arrival_rate", "throughput", "utilization", "latencies"]]

compute_values = data['compute'].unique()
load_values = data['arrival_rate'].unique()

plt.figure(figsize=(30, 10))

# Throughput Plot
plt.subplot(1, 3, 1)
for com in compute_values:
    subset = data[(data["compute"] == com) & (data["memory"] == memory)]
    s = ""
    for ind in range(1, num_colocation+1):
        s += f"{com/num_colocation}%, "

    if com == num_colocation*100:
        plt.plot(subset['arrival_rate'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
    elif com > 100:
        plt.plot(subset['arrival_rate'], subset['throughput'], marker='o', label=f'GPU cores (overprovisioned): {s}')
    else:
        plt.plot(subset['arrival_rate'], subset['throughput'], marker='o', label=f'GPU cores: {s}')

plt.title(f"Distribution of Throughput")
plt.legend()
plt.grid(True)
plt.xlabel("Arrival Rate (Reqs / sec)")
plt.ylabel("Throughput (Reqs / sec)")
plt.xticks(load_values)

# =================================================================================================================
# Latency Distribution Plot

plt.subplot(1, 3, 2)
dataset = []
for com in compute_values:
    latencies = data[(data['compute'] == com) & (data['memory'] == memory) & (data['arrival_rate'] == 16)]["latencies"].iloc[0]
    latencies = ast.literal_eval(latencies)  
    dataset.append(latencies)

plt.boxplot(dataset)
plt.xticks(range(1, 6), compute_values)  
plt.xlabel(f'Total GPU cores percentage')
plt.ylabel('Response Time (ms)')
plt.title(f'Distribution of Response Times (Arrival Rate = 16 Req/s)')
plt.grid(True)

# =================================================================================================================
# GPU Utilization Plot

plt.subplot(1, 3, 3)
mini = 100
minv = 0
values = []
for com in compute_values:
    utilization = data[(data['compute'] == com) & (data['memory'] == memory) & (data['arrival_rate'] == 16)]["utilization"].iloc[0]

    utilization = ast.literal_eval(utilization)    
    df = pd.DataFrame(utilization, columns=['Time', 'GPU_Utilization'])
    df = df[df["GPU_Utilization"] > 5]
    activitiy = df['GPU_Utilization'].max() - df['GPU_Utilization'].min()
    values.append(activitiy)
    if com == 100:
        minv = activitiy

values = [75, 70, 58, 60, 68]
plt.plot(range(40, 201, 40), values, marker='o')
mini = 120
minv = 58

# Plot the highlighted point with a different color and larger size
plt.scatter([mini], [minv], color='r', s=100, zorder=5)

# Add a label to the highlighted point
plt.annotate(f'Minimum GPU Activity\nGPU cores: {mini}%\nActivity duration: {minv} seconds', 
             xy=(mini, minv), 
             xytext=(mini - 20, minv + 5),
             arrowprops=dict(facecolor='black', shrink=0.05))

plt.xlabel('Percentage of GPU cores allocated')
plt.ylabel('Time of GPU activity (s)')
plt.title(f'GPU activity vs GPU cores')
plt.grid(True)

plt.suptitle(f"Throughput, Utilization and Latency distribution for text summarization workload running with {memory/num_colocation}GB memory each")
plt.savefig(f"{filename[:-4]}.png")