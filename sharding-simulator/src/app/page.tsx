'use client';

import { useState } from 'react';
import SimulationForm from './components/SimulationForm';
import SimulationResults from './components/SimulationResults';
import type { SimulationConfig, SimulationResults as Results } from './types';
import { Card } from '@tremor/react';

export default function Home() {
  const [results, setResults] = useState<Results | null>(null);
  const [loading, setLoading] = useState(false);

  const runSimulation = async (config: SimulationConfig) => {
    try {
      setLoading(true);
      
      const response = await fetch('http://localhost:8080/simulate-with-config', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          numNodes: config.numNodes,
          numShards: config.numShards,
          numOperators: config.numOperators,
          simulationTime: config.simulationTime,
          timeStep: config.timeStep,
          maliciousNodeRatio: config.maliciousNodeRatio,
          lotteryWinProbability: config.lotteryWinProbability,
          maliciousNodeMultiplier: config.maliciousNodeMultiplier,
          blockProductionInterval: config.blockProductionInterval,
          transactionsPerBlock: config.transactionsPerBlock,
          blockSize: config.blockSize,
          blockHeaderSize: config.blockHeaderSize,
          erHeaderSize: config.erHeaderSize,
          erBodySize: config.erBodySize,
          networkBandwidth: config.networkBandwidth,
          minNetworkDelayMean: config.minNetworkDelayMean,
          maxNetworkDelayMean: config.maxNetworkDelayMean,
          minNetworkDelayStd: config.minNetworkDelayStd,
          maxNetworkDelayStd: config.maxNetworkDelayStd,
          minGossipFanout: config.minGossipFanout,
          maxGossipFanout: config.maxGossipFanout,
          maxP2PConnections: config.maxP2PConnections,
          timeOut: config.timeOut,
          numBlocksToDownload: config.numBlocksToDownload,
          attackStartTime: config.attackStartTime,
          attackEndTime: config.attackEndTime
        })
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error('Error fetching simulation results:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-8 text-center">
          Sharding Simulator
        </h1>
        
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <Card className="p-6 bg-white dark:bg-gray-800 shadow-lg">
            <SimulationForm onSubmit={runSimulation} loading={loading} />
          </Card>
          
          {results && (
            <Card className="p-6 bg-white dark:bg-gray-800 shadow-lg">
              <SimulationResults results={results} />
            </Card>
          )}
        </div>
      </div>
    </main>
  );
}
