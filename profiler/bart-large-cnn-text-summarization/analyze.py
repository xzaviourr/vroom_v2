import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

# Read the CSV file into a DataFrame
df = pd.read_csv("bart-large-cnn-text-summarization-results.csv")
df['latency'] = df['latency']/1000
df['startup_time'] = df['startup_time']/1000

# Pivot the DataFrame to reshape it for heatmap

load_mapping = {1:"low", 5:"moderate", 20:"high"}
for load in [1, 5, 20]:
    df_load = df[df["load"] == load]
    df_pivot = df_load.pivot_table(index='memory', columns='compute', values=['startup_time', 'throughput', 'latency'])

    # Plotting the heatmap
    plt.figure(figsize=(25, 10))
    index = 1
    for variable in ['startup_time', 'throughput', 'latency']:
        plt.subplot(1, 3, index)
        index += 1
        sns.heatmap(df_pivot[variable], annot=True, cmap="YlGnBu", fmt=".2f")
        plt.title(f'Heatmap of {variable.capitalize()}')
        plt.xlabel('Compute')
        plt.ylabel('Memory')
        plt.xticks(rotation=45)
        plt.yticks(rotation=0)
        plt.gca().invert_yaxis()

    plt.suptitle(f'Bart-large-cnn-text_summarization - Load : {load_mapping[load]}', fontsize=16)  # Master title
    plt.tight_layout()
    plt.savefig(f"{load}.png")
