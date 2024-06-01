import pandas as pd
import matplotlib.pyplot as plt

single = pd.read_csv("../../bart-large-cnn-samsum-text-summarization/results.csv")

colocated2 = pd.read_csv("../samsum-2-colocation/results.csv")
colocated2["memory"] = colocated2["memory1"] + colocated2["memory2"]
colocated2["compute"] = colocated2["compute1"] + colocated2["compute2"]

colocated3 = pd.read_csv("../samsum-3-colocation/results.csv")
colocated3["memory"] = colocated3["memory1"] + colocated3["memory2"] + colocated3["memory3"]
colocated3["compute"] = colocated3["compute1"] + colocated3["compute2"] + colocated3["compute3"]

colocated4 = pd.read_csv("../samsum-4-colocation/results.csv")
colocated4["memory"] = colocated4["memory1"] + colocated4["memory2"] + colocated4["memory3"] + colocated4["memory4"]
colocated4["compute"] = colocated4["compute1"] + colocated4["compute2"] + colocated4["compute3"] + colocated4["compute4"]

single = single[['memory', 'compute', 'load', 'throughput', 'latency', 'startup_time']]
colocated2 = colocated2[['memory', 'compute', 'load', 'throughput', 'latency', 'startup_time']]
colocated3 = colocated3[['memory', 'compute', 'load', 'throughput', 'latency', 'startup_time']]
colocated4 = colocated4[['memory', 'compute', 'load', 'throughput', 'latency', 'startup_time']]

compute_power = [20, 40, 60, 80, 100]
memory = [4, 6, 8, 10, 12, 14]
for com in compute_power:
    plt.figure(figsize=(30,15))
    index = 1
    for mem in memory:
        plt.subplot(3, 2, index)
        index += 1

        single_subset = single[(single["compute"] == com) & (single["memory"] == mem)]
        colocated2_subset = colocated2[(colocated2["compute"] == com) & (colocated2["memory"] == mem)]
        colocated3_subset = colocated3[(colocated3["compute"] == com) & (colocated3["memory"] == mem)]
        colocated4_subset = colocated4[(colocated4["compute"] == com) & (colocated4["memory"] == mem)]

        plt.plot(single_subset['load'], single_subset['throughput'], marker='o', label=f'Single Workload')
        plt.plot(colocated2_subset['load'], colocated2_subset['throughput'], marker='x', label=f'Two Colocated Workloads')
        plt.plot(colocated3_subset['load'], colocated3_subset['throughput'], marker='v', label=f'Three Colocated Workloads')
        plt.plot(colocated4_subset['load'], colocated4_subset['throughput'], marker='*', label=f'Four Colocated Workloads')

        plt.title(f"Memory : {mem} GB")
        plt.legend()
        plt.grid(True)
        plt.xlabel("Load (Requests queued at once)")
        plt.ylabel("Throughput (Reqs / sec)")
        plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.suptitle(f"Throughput vs Load for different values of GPU memory with fixed GPU compute = {com}% cores")
    plt.savefig(f"ThroughputCompute{com}.png")

for com in compute_power:
    plt.figure(figsize=(30,15))
    index = 1
    for mem in memory:
        plt.subplot(3, 2, index)
        index += 1

        single_subset = single[(single["compute"] == com) & (single["memory"] == mem)]
        colocated2_subset = colocated2[(colocated2["compute"] == com) & (colocated2["memory"] == mem)]
        colocated3_subset = colocated3[(colocated3["compute"] == com) & (colocated3["memory"] == mem)]
        colocated4_subset = colocated4[(colocated4["compute"] == com) & (colocated4["memory"] == mem)]

        plt.plot(single_subset['load'], single_subset['latency'], marker='o', label=f'Single Workload')
        plt.plot(colocated2_subset['load'], colocated2_subset['latency'], marker='x', label=f'Two Colocated Workloads')
        plt.plot(colocated3_subset['load'], colocated3_subset['latency'], marker='v', label=f'Three Colocated Workloads')
        plt.plot(colocated4_subset['load'], colocated4_subset['latency'], marker='*', label=f'Four Colocated Workloads')

        plt.title(f"Memory : {mem} GB")
        plt.legend()
        plt.grid(True)
        plt.xlabel("Load (Requests queued at once)")
        plt.ylabel("Average latency per request (in sec)")
        plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.suptitle(f"Throughput vs Average Latency for different values of GPU memory with fixed GPU compute = {com}% cores")
    plt.savefig(f"LatencyCompute{com}.png")
