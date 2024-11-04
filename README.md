## Overview

The **Sharding Simulation** is a full-stack application that models a sharded blockchain network to analyze block production, shard assignments, and the impact of malicious nodes. It combines a Go backend for running complex simulations with a Next.js frontend for parameter configuration and result visualization.

## Features

- **Interactive Web Interface**: Configure and run simulations through an intuitive web UI
- **Dynamic Shard Assignment**: Nodes are assigned to shards based on lottery outcomes
- **Periodic Block Production**: Each shard produces blocks at configurable intervals
- **Malicious Nodes**: A configurable percentage of nodes behave maliciously
- **Real-time Visualization**: View simulation results through interactive charts and metrics
- **Comprehensive Analytics**: Track throughput, latency, and shard-specific statistics

## Architecture

### Backend (Go)
- Event-driven simulation engine
- RESTful API endpoints for simulation control
- Metrics collection and analysis
- Configurable simulation parameters

### Frontend (Next.js)
- React-based user interface
- Real-time metric visualization
- Parameter configuration forms
- Responsive design with Chakra UI

## Getting Started

1. **Prerequisites**
   - Go 1.22 or higher
   - Node.js 18 or higher
   - Python 3.x (for additional plotting features)

2. **Installation**

   Clone the repository and install dependencies:
   ```bash
   # Install Go dependencies
   go mod download

   # Install frontend dependencies
   cd sharding-simulator
   npm install
   ```

3. **Running the Application**

   Using PM2 (recommended):
   ```bash
   npm install -g pm2
   pm2 start ecosystem.config.js
   ```

   Or manually:
   ```bash
   # Terminal 1 - Start the backend
   go run .

   # Terminal 2 - Start the frontend
   cd sharding-simulator
   npm run dev
   ```

4. **Access the Application**
   
   Open [http://localhost:3000](http://localhost:3000) in your browser to access the web interface.

## Configuration

The simulation can be configured through the web interface or by modifying the following files:

- Frontend configuration: `sharding-simulator/.env`
  ```bash
  # API endpoint for the Go backend
  NEXT_PUBLIC_API_URL=http://localhost:8080
  ```
- Backend configuration: `config/config.go`

Key parameters include:

| Parameter | Description |
|-----------|-------------|
| Network Size | Total number of participating nodes in the blockchain network |
| Shards | Number of parallel chains the network is divided into |
| Block Interval | Time interval between consecutive block productions in each shard |
| Malicious Ratio | Proportion of nodes exhibiting malicious behavior in the network |
| Network Delays | Network latency ranges affecting message propagation between nodes |
| Operators | Number of distinct entities running validator nodes |
| Block Size | Size of each block including transactions and metadata |
| Transactions Per Block | Maximum number of transactions that can be included in a block |
| Lottery Win Probability | Chance for a node to win block production rights in a shard |
| Gossip Fanout | Number of peers each node forwards messages to during gossip |
| P2P Connections | Maximum number of peer-to-peer connections per node |
| Network Bandwidth | Available bandwidth for network communication |
| Download Timeout | Maximum time allowed for block download operations |

## Metrics and Analysis

The simulation provides real-time metrics including:

- Transaction throughput (TPS)
- Block production rates per shard
- Network latency measurements
- Malicious vs honest block ratios
- Shard-specific statistics