import pandas as pd
import matplotlib.pyplot as plt

colocated = pd.read_csv("results.csv")
colocated["memory"] = colocated["memory1"] + colocated["memory2"]
colocated["compute"] = colocated["compute1"] + colocated["compute2"]

single = pd.read_csv("../../bart-large-cnn-samsum-text-summarization/results.csv")

joint = single.copy()
joint = joint[joint["memory"] != 2]
joint = joint[joint["memory"] != 16]
joint = joint[joint["compute"] != 10]
joint = joint[joint["load"] != 1]

df1 = joint[['memory', 'compute', 'load', 'throughput', 'latency', 'startup_time']]
df2 = colocated[['memory', 'compute', 'load', 'throughput', 'latency', 'startup_time']]
merged = pd.merge(df1, df2, on=['memory', 'compute', 'load'], suffixes=('_single', '_colocated'))
print(merged)

compute_power = [20, 40, 60, 80, 100]
memory = [4, 6, 8, 10, 12, 14]
for com in compute_power:
    subset = merged[merged['compute'] == com]
    plt.figure(figsize=(25,12))
    index = 1
    for mem in memory:
        plt.subplot(2, 3, index)
        index += 1
        memsub = subset[subset['memory'] == mem]
        plt.plot(memsub['load'], memsub['throughput_single'], marker='o', label=f'Single Workload')
        plt.plot(memsub['load'], memsub['throughput_colocated'], marker='x', label=f'Two Colocated Workloads')
        plt.title(f"Memory : {mem} GB")
        plt.legend()
        plt.grid(True)
        plt.xlabel("Load (Requests queued at once)")
        plt.ylabel("Throughput (Reqs / sec)")
        plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.suptitle(f"Throughput vs Load for different values of GPU memory with fixed GPU compute = {com}% cores")
    plt.savefig(f"ThroughputCompute{com}.png")

for com in compute_power:
    subset = merged[merged['compute'] == com]
    plt.figure(figsize=(25,12))
    index = 1
    for mem in memory:
        plt.subplot(2, 3, index)
        index += 1
        memsub = subset[subset['memory'] == mem]
        plt.plot(memsub['load'], memsub['latency_single'], marker='o', label=f'Single Workload')
        plt.plot(memsub['load'], memsub['latency_colocated'], marker='x', label=f'Two Colocated Workloads')
        plt.title(f"Memory : {mem} GB")
        plt.legend()
        plt.grid(True)
        plt.xlabel("Load (Requests queued at once)")
        plt.ylabel("Average latency per request (in sec)")
        plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.suptitle(f"Throughput vs Average Latency for different values of GPU memory with fixed GPU compute = {com}% cores")
    plt.savefig(f"LatencyCompute{com}.png")

