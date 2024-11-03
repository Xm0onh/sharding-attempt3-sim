export interface SimulationConfig {
    numNodes: number;
    numOperators: number;
    numShards: number;
    simulationTime: number;
    timeStep: number;
    attackStartTime: number;
    attackEndTime: number;
    blockProductionInterval: number;
    transactionsPerBlock: number;
    maliciousNodeRatio: number;
    lotteryWinProbability: number;
    maliciousNodeMultiplier: number;
    blockSize: number;
    blockHeaderSize: number;
    erHeaderSize: number;
    erBodySize: number;
    networkBandwidth: number;
    minNetworkDelayMean: number;
    maxNetworkDelayMean: number;
    minNetworkDelayStd: number;
    maxNetworkDelayStd: number;
    minGossipFanout: number;
    maxGossipFanout: number;
    maxP2PConnections: number;
    timeOut: number;
    numBlocksToDownload: number;
}
  
export interface SimulationResults {
    performance: {
        transactions_per_second: number;
    };
    transaction_size_bytes: number;
    block_size_kb: number;
    block_production: {
        [key: string]: {
            honest_blocks: number;
            malicious_blocks: number;
            total_blocks: number;
        };
    };
    network_metrics: {
        block_header_delay_ms: number;
        block_broadcast_delays_ms: {
            [key: string]: number;
        };
        block_download_delays_ms: {
            [key: string]: number;
        };
    };
}