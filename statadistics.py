import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import sys

CSV_FILE = 'benchmark_data.csv'
OUTPUT_IMAGE = 'benchmark_analysis.png'
COLORS = ["#4CAF50", "#F44336"]

def load_and_process_data(filepath):
    """Loads CSV and adds necessary calculated columns."""
    try:
        df = pd.read_csv(filepath)
    except FileNotFoundError:
        print(f"Error: {filepath} not found. Run the Go server first.")
        sys.exit(1)

    # Convert Nanoseconds to Microseconds
    df['duration_us'] = df['duration_ns'] / 1000 
    
    df['request_id'] = range(1, len(df) + 1)
    
    return df

def plot_latency_distribution(df, axes):
    """Draws a Box Plot showing the spread of response times."""
    sns.boxplot(
        x='endpoint', 
        y='duration_us', 
        data=df, 
        ax=axes, 
        hue='endpoint', 
        palette=COLORS, 
        legend=False
    )
    axes.set_title('Summary: Latency Distribution', fontsize=12, fontweight='bold')
    axes.set_ylabel('Duration (µs)')
    axes.set_xlabel('')

def plot_network_usage(df, axes):
    """Draws a Bar Chart showing the payload size per endpoint."""
    sns.barplot(
        x='endpoint', 
        y='size_bytes', 
        data=df, 
        ax=axes, 
        hue='endpoint', 
        palette=COLORS,
        errorbar=None,
        legend=False
    )
    axes.set_title('Summary: Network Payload Size', fontsize=12, fontweight='bold')
    axes.set_ylabel('Size (Bytes)')
    axes.set_xlabel('')
    
    # Add text labels on top of bars
    for container in axes.containers:
        axes.bar_label(container, padding=3)

# For showing all requests
def plot_request_timeline(df, axes):
    """Draws a Scatter Plot showing every single request over time."""

    sns.scatterplot(
        data=df,
        x='request_id',
        y='duration_us',
        hue='endpoint',
        #style='endpoint',
        palette=COLORS,
        s=60,
        alpha=0.7,
        ax=axes
    )
    axes.set_title('Timeline: Every Request Captured', fontsize=12, fontweight='bold')
    axes.set_ylabel('Duration (µs)')
    axes.set_xlabel('Request')
    axes.legend(title="Endpoint", loc='upper right')

def main():
    df = load_and_process_data(CSV_FILE)

    # Theme for canvas 
    sns.set_theme(style="whitegrid")
    fig = plt.figure(figsize=(16, 12))
    
    # Grid: 2 rows, 2 columns.
    # Top row takes 1 unit height, Bottom row takes 1.2 units height.
    grid = plt.GridSpec(2, 2, height_ratios=[1, 1.2], hspace=0.3)

    # Top Left
    ax1 = fig.add_subplot(grid[0, 0])
    plot_latency_distribution(df, ax1)

    # Top Right
    ax2 = fig.add_subplot(grid[0, 1])
    plot_network_usage(df, ax2)

    # Bottom : works like flex/expanded in HTML/flutter (kind of)
    ax3 = fig.add_subplot(grid[1, :])
    plot_request_timeline(df, ax3)
    
    plt.suptitle('API Performance Benchmark: Optimal vs Double decode', fontsize=16, y=0.95)
    plt.savefig(OUTPUT_IMAGE)
    print(f"Success! Chart saved to: {OUTPUT_IMAGE}")
    plt.show()

if __name__ == "__main__":
    main()