import pandas as pd
import matplotlib.pyplot as plt

# Load the data
data = pd.read_csv("log_entries.csv")

# Convert timestamp columns to datetime data type
data['registrationts'] = pd.to_datetime(data['registrationts'])
data['responsets'] = pd.to_datetime(data['responsets'])

# Flatten the timestamp range to every second
timestamps = pd.date_range(data['registrationts'].min(), data['responsets'].max(), freq='S')

instances = data["instance"].unique()

index = 1
mapping = {1:"blue", 2:"green",3:"yellow",4:"red",5:"orange"}
labels = {1:"16 GB, 100%", 2:"16 GB, 100%", 3:"2 GB, 40%", 4: "4 GB, 20%"}

plt.figure(figsize=(10,10))
final = []
for instance in instances:
    plt.subplot(3, 1, index)
    # Count the occurrences of each timestamp
    counts = pd.Series(index=timestamps, data=0)
    for idx, row in data[data["instance"] == instance].iterrows():
        counts.loc[row['responsets']] += 1

    final.append(counts.values)
    # Plot the line chart
    plt.ylabel('Throughput')
    plt.plot(range(len(counts.index)), counts.values, color=mapping[index], label=f"{labels[index]}")
    plt.grid(True)
    plt.legend()

    index += 1
    if index == 5:
        break
    
max_length = 0
for i in range(len(final)):
    max_length = max(max_length, len(final[i]))

total_throughput = [0 for x in range(max_length)]
for i in range(max_length):
    for j in range(len(final)):
        try:
            total_throughput[i] += final[j][i]
        except:
            pass

plt.subplot(3, 1, 3)
plt.plot(range(max_length), total_throughput, color='black', label=f"Total Throughput")

# plt.xticks(range(len(counts.index)))
plt.xlabel('Timeline')
plt.ylabel('Throughput')
plt.suptitle('Throughput timeline for the workload with constant 3 reqs/sec arrival rate.')
plt.grid(True)
plt.legend()
plt.savefig("simulation.png")
plt.show()