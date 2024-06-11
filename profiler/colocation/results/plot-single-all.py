import pandas as pd
import matplotlib.pyplot as plt

single = pd.read_csv("cnn-1-fixed-memory-all.csv")
single['memory'] = single['memory1']
single['compute'] = single['compute1']
single = single[['memory', 'compute', 'arrival_rate', 'throughput']]

compute_power = [20, 40, 60, 80, 100]
memory = [2, 4, 6, 8, 10, 12, 14, 16]

plt.figure(figsize=(20, 40))
index = 1
for mem in memory:
    plt.subplot(4, 2, index)
    index += 1
    for com in compute_power:
        single_subset = single[(single["compute"] == com) & (single['memory'] == mem)]
        plt.plot(single_subset['arrival_rate'], single_subset['throughput'], marker='o', label=f"GPU cores: {com}%")
    plt.title(f"GPU Memory: {mem} GB")
    plt.legend()
    plt.grid(True)
    plt.xlabel("Load (Requests queued at once)")
    plt.ylabel("Throughput (Reqs / sec)")
    plt.xticks([2, 4, 6, 8, 10, 12, 14, 16])

plt.suptitle(f"Profile of large CNN for text summarization.")
plt.savefig(f"single-profile.png")