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

colocated2o = pd.read_csv("../samsum-2-overprovisioned/results.csv")
colocated2o["memory"] = colocated2o["memory1"] + colocated2o["memory2"]
colocated2o["compute"] = colocated2o["compute1"] + colocated2o["compute2"]

colocated3o = pd.read_csv("../samsum-3-overprovisioned/results.csv")
colocated3o["memory"] = colocated3o["memory1"] + colocated3o["memory2"] + colocated3o["memory3"]
colocated3o["compute"] = colocated3o["compute1"] + colocated3o["compute2"] + colocated3o["compute3"]

colocated4o = pd.read_csv("../samsum-4-overprovisioned/results.csv")
colocated4o["memory"] = colocated4o["memory1"] + colocated4o["memory2"] + colocated4o["memory3"] + colocated4o["memory4"]
colocated4o["compute"] = colocated4o["compute1"] + colocated4o["compute2"] + colocated4o["compute3"] + colocated4o["compute4"]

# =========================================================================================

colocated2 = colocated2[["memory", "compute", "load", "latency", "throughput"]]
colocated2o = colocated2o[["memory", "compute", "load", "latency", "throughput"]]
dataset2 = pd.concat([colocated2, colocated2o], ignore_index=True)

colocated3 = colocated3[["memory", "compute", "load", "latency", "throughput"]]
colocated3o = colocated3o[["memory", "compute", "load", "latency", "throughput"]]
dataset3 = pd.concat([colocated3, colocated3o], ignore_index=True)

colocated4 = colocated4[["memory", "compute", "load", "latency", "throughput"]]
colocated4o = colocated4o[["memory", "compute", "load", "latency", "throughput"]]
dataset4 = pd.concat([colocated4, colocated4o], ignore_index=True)

memory = list(dataset2["memory"].unique())

for mem in memory:
    plt.figure(figsize=(28,20))
    
    plt.subplot(2, 2, 1)
    compute = list(single["compute"].unique())
    for com in compute:
        subset = single[(single["compute"] == com) & (single["memory"] == mem)]
        if com == 100:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
        else:
            plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores: {com}%')
    
    plt.title(f"Single workload with GPU Memory: {mem} GB")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.subplot(2, 2, 2)
    compute = list(dataset2["compute"].unique())
    for com in compute:
        subset = dataset2[(dataset2["compute"] == com) & (dataset2["memory"] == mem)]
        if com == 200:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
        elif com > 100:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'GPU cores (overprovisioned): {com/2}% and {com/2}%')
        else:
            plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores: {com/2}% and {com/2}%')
    
    plt.title(f"Two Colocated workloads with equal Memory : {mem/2} GB and {mem/2} GB")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.subplot(2, 2, 3)
    compute = list(dataset3["compute"].unique())
    for com in compute:
        subset = dataset3[(dataset3["compute"] == com) & (dataset3["memory"] == mem)]
        if com == 300:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
        elif com > 100:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'GPU cores (overprovisioned): {round(com/3, 2)}%, {round(com/3, 2)}%, and {round(com/3, 2)}%')
        else:
            plt.plot(subset['load'], subset['throughput'], marker='o', label=f'GPU cores = {round(com/3, 2)}%, {round(com/3, 2)}%, and {round(com/3, 2)}%')
        
    
    plt.title(f"Three Colocated workloads with equal Memory : {mem/3} GB, {mem/3} GB, and {mem/3} GB")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.subplot(2, 2, 4)
    compute = list(dataset4["compute"].unique())
    for com in compute:
        subset = dataset4[(dataset4["compute"] == com) & (dataset4["memory"] == mem)]
        if com == 400:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'Default no limit setup')
        elif com > 100:
            plt.plot(subset['load'], subset['throughput'], marker='X', markersize=12, linewidth = 3, label=f'GPU cores (overprovisioned): {com/4}%, {com/4}%, {com/4}%, and {com/4}%')
        else:
            plt.plot(subset['load'], subset['throughput'], marker='o', label=f"GPU cores = {com/4}%, {com/4}%, {com/4}%, and {com/4}%")
    
    plt.title(f"Four Colocated workloads with equal Memory : {mem/4} GB, {mem/4} GB, {mem/4} GB, and {mem/4} GB")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 8, 16, 32, 64, 128])

    plt.suptitle(f"These graphs represents how the throughput varies by varying the overall GPU cores for colocated workloads.\nGPU memory is set to {mem} GB, and compute limit is divided equally among the workloads.\nAt different levels of colocation (4 different graphs), different values of GPU cores gives different performance as shown.\n It can be seen that maximum over-provisioning or the no limit case (which is currently available in the market) does not necessarily gives the best performance.")
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
