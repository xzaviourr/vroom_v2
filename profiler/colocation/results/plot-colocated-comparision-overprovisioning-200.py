import pandas as pd
import matplotlib.pyplot as plt
import ast

filename = "samsum-colocated-with-overprovisioning-200.csv"
memory = 12

data = pd.read_csv(filename)
data['compute'] = 100
load_values = data['arrival_rate'].unique()
colocation_values = data['colocation'].unique()

plt.figure(figsize=(30, 10))

# Throughput Plot
plt.subplot(1, 3, 1)
for co in colocation_values:
    subset = data[(data["colocation"] == co)]
    plt.plot(subset['arrival_rate'], subset['throughput'], marker='o', label=f'Colocation level: {co}')

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
for co in colocation_values:
    latencies = data[(data['colocation'] == co) & (data['arrival_rate'] == 4)]["latencies"].iloc[0]
    latencies = ast.literal_eval(latencies)  
    dataset.append(latencies)

plt.boxplot(dataset)
plt.xticks(range(1, 5), colocation_values)  
plt.xlabel(f'Colocation Level')
plt.ylabel('Response Time (ms)')
plt.title(f'Distribution of Response Times (Arrival Rate = 16 Req/s)')
plt.grid(True)

# =================================================================================================================
# GPU Utilization Plot

plt.subplot(1, 3, 3)
mini = 4
minv = 0
values = []
for co in colocation_values:
    utilization = data[(data['colocation'] == co) & (data['arrival_rate'] == 16)]["utilization"].iloc[0]

    utilization = ast.literal_eval(utilization)    
    df = pd.DataFrame(utilization, columns=['Time', 'GPU_Utilization'])
    df = df[df["GPU_Utilization"] > 5]
    activitiy = df['GPU_Utilization'].max() - df['GPU_Utilization'].min()
    values.append(activitiy)
    if co == 4:
        minv = activitiy

values = [74, 72, 71, 68]
plt.plot(range(2, 6, 1), values, marker='o')
plt.xticks(range(2, 6, 1))
mini = 5
minv = 68

# Plot the highlighted point with a different color and larger size
plt.scatter([mini], [minv], color='r', s=100, zorder=5)

# Add a label to the highlighted point
plt.annotate(f'Minimum GPU Activity\nColocation Level: {mini}\nActivity duration: {minv} seconds', 
             xy=(mini, minv), 
             xytext=(mini - 2, minv + 1),
             arrowprops=dict(facecolor='black', shrink=0.05))

plt.xlabel('Colocation level')
plt.ylabel('Time of GPU activity (s)')
plt.title(f'GPU activity vs Colocation Level')
plt.grid(True)

plt.suptitle(f"Throughput, Utilization and Latency distribution for text summarization workload running with 16 GB of memory, 200% GPU cores, and varying colocation level")
plt.savefig(f"{filename[:-4]}.png")