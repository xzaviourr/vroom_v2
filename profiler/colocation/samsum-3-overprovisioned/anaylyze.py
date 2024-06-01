import pandas as pd
import matplotlib.pyplot as plt

single = pd.read_csv("../../bart-large-cnn-samsum-text-summarization/results.csv")

colocated2 = pd.read_csv("../samsum-2-overprovisioned/results.csv")
colocated2["memory"] = colocated2["memory1"] + colocated2["memory2"]
colocated2["compute"] = colocated2["compute1"] + colocated2["compute2"]

colocated3 = pd.read_csv("results.csv")
colocated3["memory"] = colocated3["memory1"] + colocated3["memory2"] + colocated3["memory3"]
colocated3["compute"] = colocated3["compute1"] + colocated3["compute2"] + colocated3["compute3"]

colocated4 = pd.read_csv("../samsum-4-overprovisioned/results.csv")
colocated4["memory"] = colocated4["memory1"] + colocated4["memory2"] + colocated4["memory3"] + colocated4["memory4"]
colocated4["compute"] = colocated4["compute1"] + colocated4["compute2"] + colocated4["compute3"] + colocated4["compute4"]

memory = list(colocated3["memory"].unique())

for mem in memory:
    plt.figure(figsize=(25,15))
    
    plt.subplot(2, 2, 1)
    compute = list(single["compute"].unique())
    for com in compute:
        subset = single[(single["compute"] == com) & (single["memory"] == mem)]
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'Compute = {com}%')
    
    plt.title(f"Single workload with Memory : {mem} GB")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.subplot(2, 2, 2)
    compute = list(colocated2["compute"].unique())
    for com in compute:
        subset = colocated2[(colocated2["compute"] == com) & (colocated2["memory"] == mem)]
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'Compute = {com}%')
    
    plt.title(f"Two Colocated workloads with equal Memory : {mem} GB and equal Compute cores")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.subplot(2, 2, 3)
    compute = list(colocated3["compute"].unique())
    for com in compute:
        subset = colocated3[(colocated3["compute"] == com) & (colocated3["memory"] == mem)]
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'Compute = {com}%')
    
    plt.title(f"Three Colocated workloads with equal Memory : {mem} GB and equal Compute cores")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.subplot(2, 2, 4)
    compute = list(colocated4["compute"].unique())
    for com in compute:
        subset = colocated4[(colocated4["compute"] == com) & (colocated4["memory"] == mem)]
        plt.plot(subset['load'], subset['throughput'], marker='o', label=f'Compute = {com}%')
    
    plt.title(f"Four Colocated workloads with equal Memory : {mem} GB and equal Compute cores")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.suptitle(f"These graphs represents how the throughput varies by varying the level of over-provisioning for colocated workloads.\nGPU memory is set to {mem} GB, and compute limit is divided equally among the workloads.\nAt different levels of colocation (4 different graphs), different levels of overprovisioning gives varied performance as shown.\n It can be seen that maximum over-provisioning (which is available currently in the market) does not necessarily gives the best performance.")
    plt.savefig(f"ThroughputMemory{mem}.png")

# for com in compute_power:
#     plt.figure(figsize=(30,15))
#     index = 1
#     for mem in memory:
#         plt.subplot(3, 2, index)
#         index += 1

#         single_subset = single[(single["compute"] == com) & (single["memory"] == mem)]
#         colocated2_subset = colocated2[(colocated2["compute"] == com) & (colocated2["memory"] == mem)]
#         colocated3_subset = colocated3[(colocated3["compute"] == com) & (colocated3["memory"] == mem)]
#         colocated4_subset = colocated4[(colocated4["compute"] == com) & (colocated4["memory"] == mem)]

#         plt.plot(single_subset['load'], single_subset['latency'], marker='o', label=f'Single Workload')
#         plt.plot(colocated2_subset['load'], colocated2_subset['latency'], marker='x', label=f'Two Colocated Workloads')
#         plt.plot(colocated3_subset['load'], colocated3_subset['latency'], marker='v', label=f'Three Colocated Workloads')
#         plt.plot(colocated4_subset['load'], colocated4_subset['latency'], marker='*', label=f'Four Colocated Workloads')

#         plt.title(f"Memory : {mem} GB")
#         plt.legend()
#         plt.grid(True)
#         plt.xlabel("Load (Requests queued at once)")
#         plt.ylabel("Average latency per request (in sec)")
#         plt.xticks([2, 4, 8, 16, 32, 64, 128])

#     plt.suptitle(f"Throughput vs Average Latency for different values of GPU memory with fixed GPU compute = {com}% cores")
#     plt.savefig(f"LatencyCompute{com}.png")
