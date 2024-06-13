import matplotlib.pyplot as plt

cpu = [0.55, 0.55, 0.55, 0.55, 0.55]
gpu = [0.78, 0.84, 0.97, 1.2, 1.35]
cores = [20, 40, 60, 80, 100]

plt.figure(figsize=(6, 6))
plt.plot(cores, cpu, marker='o', label="cpu")
plt.plot(cores, gpu, marker='o', label="gpu")
plt.xticks(cores)
plt.yticks([0, 0.2, 0.4, 0.6, 0.8, 1, 1.2, 1.4])
plt.legend()
plt.xlabel("GPU cores (in percentage)")
plt.ylabel("Throughput (reqs/sec)")
plt.title("Text summarization function running on CPU vs GPU")
plt.savefig("cpu_vs_gpu.png")