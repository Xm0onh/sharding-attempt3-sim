import matplotlib.pyplot as plt

def parse_metrics_report(file_path):
    timestamps = []
    tps = []
    avg_latency = []
    malicious_rotations = []

    with open(file_path, 'r') as file:
        lines = file.readlines()
        for line in lines:
            if line.startswith("Time:"):
                parts = line.split(',')
                timestamp = int(parts[0].split(':')[1].strip())
                tps_value = float(parts[3].split(':')[1].strip())
                latency_value = float(parts[4].split(':')[1].strip().split()[0])
                timestamps.append(timestamp)
                tps.append(tps_value)
                avg_latency.append(latency_value)
                # Ensure malicious_rotations has the same length as timestamps
                if len(malicious_rotations) < len(timestamps):
                    malicious_rotations.append(0)
            elif line.startswith("  Malicious Shard Rotations This Step:"):
                rotations = int(line.split(':')[1].strip())
                if len(malicious_rotations) == len(timestamps):
                    malicious_rotations[-1] = rotations
                else:
                    malicious_rotations.append(rotations)

    return timestamps, tps, avg_latency, malicious_rotations

def plot_metrics(timestamps, tps, avg_latency, malicious_rotations):
    fig, ax1 = plt.subplots()

    color = 'tab:blue'
    ax1.set_xlabel('Time')
    ax1.set_ylabel('TPS', color=color)
    ax1.plot(timestamps, tps, color=color, label='TPS')
    ax1.tick_params(axis='y', labelcolor=color)

    ax2 = ax1.twinx()
    color = 'tab:red'
    ax2.set_ylabel('Malicious Shard Rotations', color=color)
    ax2.plot(timestamps, malicious_rotations, color=color, label='Malicious Shard Rotations')
    ax2.tick_params(axis='y', labelcolor=color)

    fig.tight_layout()
    plt.title('TPS and Malicious Shard Rotations Over Time')
    plt.show()

    plt.figure()
    plt.plot(timestamps, avg_latency, label='Average Latency', color='tab:green')
    plt.xlabel('Time')
    plt.ylabel('Average Latency')
    plt.title('Average Latency Over Time')
    plt.legend()
    plt.show()

if __name__ == "__main__":
    file_path = 'metrics_report.txt'
    timestamps, tps, avg_latency, malicious_rotations = parse_metrics_report(file_path)
    plot_metrics(timestamps, tps, avg_latency, malicious_rotations)