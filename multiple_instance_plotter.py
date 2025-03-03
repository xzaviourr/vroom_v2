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
labels = {1:"2 GB, 20%", 2:"2 GB, 20%", 3:"2 GB, 40%", 4: "4 GB, 20%"}
print(instances)
plt.figure(figsize=(10,10))
for instance in instances:
    # Count the occurrences of each timestamp
    counts = pd.Series(index=timestamps, data=0)
    for idx, row in data[data["instance"] == instance].iterrows():
        counts.loc[row['responsets']] += 1

    # Plot the line chart
    plt.plot(range(len(counts.index)), counts.values, color=mapping[index], label=f"{labels[index]}")
    index += 1
    if index == 5:
        break

# plt.xticks(range(len(counts.index)))
plt.xlabel('Timeline')
plt.ylabel('Frequency')
plt.title('Throughput timeline for the workload with constant 3 reqs/sec arrival rate.')
plt.grid(True)
plt.legend()
plt.savefig("simulation.png")
plt.show()