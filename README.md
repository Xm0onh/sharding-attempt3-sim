## Overview

The **Sharding Simulation** models a sharded blockchain network to analyze block production, shard assignments, and the impact of malicious nodes. It enables testing the scalability and performance of a sharded architecture under various configurations.

## Features

- **Dynamic Shard Assignment**: Nodes are assigned to shards based on lottery outcomes.
- **Periodic Block Production**: Each shard produces blocks at configurable intervals.
- **Malicious Nodes**: A configurable percentage of nodes behave maliciously, influencing block creation.
- **Comprehensive Metrics**: Tracks total and malicious blocks, throughput (TPS), latency, and shard-specific statistics.
- **Scalability Testing**: Adjust parameters like the number of nodes and shards to evaluate system scalability.

## How It Works

1. **Initialization**
   - **Configuration**: Set simulation parameters in `config/config.go`.
   - **Shard and Node Setup**: Initialize shards and create nodes, initially unassigned to any shard.

2. **Event Scheduling**
   - **Lottery Events**: Nodes attempt to join shards via lottery at regular intervals.
   - **Shard Block Production Events**: Shards produce blocks every `BlockProductionInterval` time steps.
   - **Metrics Events**: Collect and record metrics at each time step.

3. **Event Processing**
   - **Lottery Event**: Nodes participate in a lottery to join shards. Winners are assigned to shards and immediately produce a block.
   - **Shard Block Production Event**: Shards select a node to produce a block and broadcast it to shard members.
   - **Message Event**: Nodes receive and process incoming blocks, updating their known blocks.
   - **Metrics Event**: Gather and log metrics such as block counts, TPS, and latency.

## Configuration

Modify the simulation parameters in `config/config.go`:

| Parameter                   | Type      | Description                                                   | Default Value |
|-----------------------------|-----------|---------------------------------------------------------------|---------------|
| `NumNodes`                  | `int`     | Total number of nodes in the network.                         | 10000         |
| `NumShards`                 | `int`     | Total number of shards in the network.                        | 10            |
| `SimulationTime`            | `int64`   | Total duration of the simulation in time units.               | 120           |
| `TimeStep`                  | `int64`   | Simulation advances in increments of this time unit.          | 1             |
| `NetworkDelayMean`          | `int64`   | Average network delay in time units for block propagation.    | 5             |
| `NetworkDelayStd`           | `int64`   | Standard deviation of network delay in time units.            | 2             |
| `MaliciousNodeRatio`        | `float64` | Percentage of nodes that behave maliciously.                  | 0.1 (10%)     |
| `LotteryWinProbability`     | `float64` | Probability of an honest node winning the lottery per attempt.| 0.001         |
| `MaliciousNodeMultiplier`   | `int`     | Additional lottery attempts for malicious nodes.              | 5             |
| `BlockProductionInterval`   | `int64`   | Time steps between each shard's block production events.      | 6             |
| `TransactionsPerBlock`      | `int`     | Number of transactions contained in each block.               | 100           |
| `AttackStartTime`           | `int64`   | Time step when the attack starts.                             | 20            |
| `AttackEndTime`             | `int64`   | Time step when the attack ends.                               | 60            |
| `AttackType`                | `AttackType` | Type of attack to simulate.                                  | GrindingAttack|
| `AttackSchedule`            | `map[int64]AttackType` | Schedule of attacks with start and end times.               | Initialized by `InitializeAttackSchedule` |

## Metrics and Analysis

The simulation collects the following metrics:

- **Total Blocks**: Cumulative number of blocks produced across all shards.
- **Malicious Blocks**: Number of blocks produced by malicious nodes.
- **Throughput (TPS)**: Transactions Per Second, calculated as Total Transactions / Simulation Time.
- **Latency**: Approximate time from transaction submission to confirmation, calculated as BlockProductionInterval + Average Network Delay.
- **Shard Statistics**:
  - Honest Nodes: Number of honest nodes in each shard.
  - Malicious Nodes: Number of malicious nodes in each shard.
  - Honest Blocks: Blocks produced by honest nodes in each shard.
  - Malicious Blocks: Blocks produced by malicious nodes in each shard.

## Event Mechanism

The simulation uses an event-driven mechanism to manage and process various activities. Events are scheduled and processed based on their timestamps to simulate the passage of time in the network.

### Event Types

- **LotteryEvent**: Nodes attempt to join shards via a lottery.
- **MessageEvent**: Nodes receive and process incoming blocks.
- **AttackEvent**: Simulates attacks on the network.
- **MetricsEvent**: Collects and logs metrics.
- **ShardBlockProductionEvent**: Shards produce blocks at regular intervals.

### Event Queue

Events are managed using a priority queue to ensure they are processed in the correct order. The event queue is implemented using a heap, which allows efficient insertion and removal of events based on their timestamps.

- **Initialization**: The event queue is initialized using `NewEventQueue()`.
- **Event Handling**: Events are added to the queue using `Push()` and processed in order using `Pop()`.
- **Empty Check**: The queue can be checked for emptiness using `IsEmpty()`.

This event-driven approach allows the simulation to accurately model the timing and sequence of activities in a sharded blockchain network.

## Running the Simulation

To run the simulation and visualize the results, follow these steps:

1. Run the simulation:
   ```
   go run .
   ```
   This will execute the simulation and generate a `metrics_report.txt` file with the simulation results.

2. Create a Python virtual environment (optional but recommended):
   ```
   python -m venv venv
   source venv/bin/activate  # On Windows, use `venv\Scripts\activate`
   ```

3. Install the required Python packages:
   ```
   pip install matplotlib pandas
   ```

4. Run the plotting script:
   ```
   python plot.py
   ```
   This script will read the `metrics_report.txt` file, process the data, and generate visualizations of the simulation results.

The generated plots will provide visual insights into various aspects of the simulation, such as transaction throughput, block production rates, and the effects of attacks on the network.


<!-- ## Extending the Simulation

Enhance the simulation by:

1. **Implementing Attack Scenarios**: 
   - Add specific attack types in `attack/attack.go` to study their effects.
   - Examples: Sybil attacks, Eclipse attacks, or Selfish mining.

2. **Detailed Transaction Modeling**: 
   - Simulate individual transactions within blocks for more granular metrics.
   - Track transaction confirmation times and success rates.

4. **Advanced Metrics**: 
   - Collect additional metrics such as:
     - Fork rates
     - Orphaned blocks
     - Node churn rates
     - Network partitioning events

5. **Visualization**: 
   - Integrate with visualization tools to graphically represent metrics over time.
   - Create dynamic charts and graphs for real-time simulation monitoring. -->

<!-- sharding/
├── main.go
├── config/
│   └── config.go
├── simulation/
│   └── simulation.go
├── event/
│   └── event.go
├── node/
│   └── node.go
├── shard/
│   └── shard.go
├── block/
│   └── block.go
├── lottery/
│   └── lottery.go
├── attack/
│   └── attack.go
├── metrics/
│   └── metrics.go
├── utils/
    ├── constants.go
    └── random.go -->