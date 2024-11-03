import matplotlib.pyplot as plt
import numpy as np
import re

def parse_metrics_report(file_path):
    metrics = {
        'block_size': 0,
        'txn_per_block': 0,
        'shard_stats': {},  # {shard_id: {'honest_blocks': [], 'malicious_blocks': []}}
        'network_metrics': {
            'broadcast_delay': {},  # {shard_id: avg_delay}
            'header_delay': 0,
            'download_delay': {},   # {shard_id: avg_delay}
        },
        'tps': 0
    }
    
    with open(file_path, 'r') as file:
        lines = file.readlines()
        
        for i, line in enumerate(lines):
            # Parse basic configuration
            if "Size of each Block in kilo bytes:" in line:
                metrics['block_size'] = int(re.findall(r'\d+', line)[0])
            elif "Number of transactions per block:" in line:
                metrics['txn_per_block'] = int(re.findall(r'\d+', line)[0])
            
            # Parse shard statistics
            elif "=== Shard" in line:
                shard_id = int(re.findall(r'Shard (\d+)', line)[0])
                if shard_id not in metrics['shard_stats']:
                    metrics['shard_stats'][shard_id] = {'honest_blocks': 0, 'malicious_blocks': 0}
                
                # Parse next lines for block counts
                malicious_line = lines[i + 1]
                honest_line = lines[i + 2]
                
                metrics['shard_stats'][shard_id]['malicious_blocks'] = int(re.findall(r'\d+', malicious_line)[-1])
                metrics['shard_stats'][shard_id]['honest_blocks'] = int(re.findall(r'\d+', honest_line)[-1])
            
            # Parse network metrics
            elif "Average Block Broadcast Delay per Shard:" in line:
                j = i + 1
                while "Shard" in lines[j]:
                    shard_id = int(re.findall(r'Shard (\d+)', lines[j])[0])
                    delay = float(re.findall(r'(\d+\.\d+)ms', lines[j])[0])
                    metrics['network_metrics']['broadcast_delay'][shard_id] = delay
                    j += 1
            
            elif "Average Block Header Delay:" in line:
                metrics['network_metrics']['header_delay'] = float(re.findall(r'(\d+\.\d+)ms', line)[0])
            
            elif "Average Block Download Delay per Shard:" in line:
                j = i + 1
                while j < len(lines) and "Shard" in lines[j]:
                    shard_id = int(re.findall(r'Shard (\d+)', lines[j])[0])
                    delay = float(re.findall(r'(\d+\.\d+)ms', lines[j])[0])
                    metrics['network_metrics']['download_delay'][shard_id] = delay
                    j += 1
            
            # Parse TPS
            elif "Transactions Per Second (TPS):" in line:
                metrics['tps'] = float(re.findall(r'(\d+\.\d+)', line)[0])
    
    return metrics

def plot_metrics(metrics):
    # 1. Block Distribution across Shards
    plt.figure(figsize=(12, 6))
    shard_ids = list(metrics['shard_stats'].keys())
    honest_blocks = [metrics['shard_stats'][sid]['honest_blocks'] for sid in shard_ids]
    malicious_blocks = [metrics['shard_stats'][sid]['malicious_blocks'] for sid in shard_ids]
    
    x = np.arange(len(shard_ids))
    width = 0.35
    
    plt.bar(x - width/2, honest_blocks, width, label='Honest Blocks')
    plt.bar(x + width/2, malicious_blocks, width, label='Malicious Blocks')
    plt.xlabel('Shard ID')
    plt.ylabel('Number of Blocks')
    plt.title('Block Distribution Across Shards')
    plt.xticks(x, [f'Shard {sid}' for sid in shard_ids])
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.show()

    # 2. Network Delays Comparison
    plt.figure(figsize=(12, 6))
    shard_ids = list(metrics['network_metrics']['broadcast_delay'].keys())
    broadcast_delays = [metrics['network_metrics']['broadcast_delay'][sid] for sid in shard_ids]
    download_delays = [metrics['network_metrics']['download_delay'][sid] for sid in shard_ids]
    
    x = np.arange(len(shard_ids))
    width = 0.35
    
    plt.bar(x - width/2, broadcast_delays, width, label='Broadcast Delay')
    plt.bar(x + width/2, download_delays, width, label='Download Delay')
    plt.axhline(y=metrics['network_metrics']['header_delay'], color='r', linestyle='--', 
                label='Average Header Delay')
    
    plt.xlabel('Shard ID')
    plt.ylabel('Delay (ms)')
    plt.title('Network Delays by Shard')
    plt.xticks(x, [f'Shard {sid}' for sid in shard_ids])
    plt.legend()
    plt.grid(True, alpha=0.3)
    plt.show()

    # 3. Block Production Efficiency
    plt.figure(figsize=(12, 6))
    for shard_id in metrics['shard_stats']:
        total_blocks = (metrics['shard_stats'][shard_id]['honest_blocks'] + 
                       metrics['shard_stats'][shard_id]['malicious_blocks'])
        honest_ratio = metrics['shard_stats'][shard_id]['honest_blocks'] / total_blocks if total_blocks > 0 else 0
        
        plt.bar(f'Shard {shard_id}', honest_ratio * 100)
    
    plt.ylabel('Honest Block Percentage (%)')
    plt.title('Honest Block Production Ratio by Shard')
    plt.grid(True, alpha=0.3)
    plt.show()

    # 4. System Performance Summary
    plt.figure(figsize=(8, 6))
    performance_metrics = {
        'TPS': metrics['tps'],
        'Avg Header\nDelay (ms)': metrics['network_metrics']['header_delay'],
        'Avg Block\nSize (KB)': metrics['block_size']
    }
    
    plt.bar(performance_metrics.keys(), performance_metrics.values())
    plt.title('System Performance Summary')
    plt.grid(True, alpha=0.3)
    plt.show()

if __name__ == "__main__":
    file_path = 'simulation_report.txt'
    metrics = parse_metrics_report(file_path)
    plot_metrics(metrics)