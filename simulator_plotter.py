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

# Count the occurrences of each timestamp
counts = pd.Series(index=timestamps, data=0)
for index, row in data.iterrows():
    counts.loc[row['responsets']] += 1

# Plot the histogram
plt.figure(figsize=(10, 6))
plt.bar(range(len(counts.index)), counts.values, width=1, color='skyblue', edgecolor='black', label="Variant 20%,2GB")

# plt.xticks(range(len(counts.index)))
plt.xlabel('Timeline')
plt.ylabel('Frequency')
plt.title('Throughput timeline for the workload with costant 0.2 reqs/sec arrival rate')
plt.grid(True)
plt.legend()
plt.savefig("simulation.png")
plt.show()
