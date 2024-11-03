import { NextResponse } from 'next/server';
import type { SimulationConfig } from '../../types';

export async function POST(request: Request) {
  try {
    const config: SimulationConfig = await request.json();

    // Call your Go simulation here
    // For now, returning mock data
    const mockResults = {
      tps: 1000 + Math.random() * 500,
      shardStats: {
        '0': {
          honestNodes: Math.floor(config.numNodes * (1 - config.maliciousNodeRatio) / config.numShards),
          maliciousNodes: Math.floor(config.numNodes * config.maliciousNodeRatio / config.numShards),
          honestBlocks: Math.floor(Math.random() * 1000),
          maliciousBlocks: Math.floor(Math.random() * 100),
        },
        '1': {
          honestNodes: Math.floor(config.numNodes * (1 - config.maliciousNodeRatio) / config.numShards),
          maliciousNodes: Math.floor(config.numNodes * config.maliciousNodeRatio / config.numShards),
          honestBlocks: Math.floor(Math.random() * 1000),
          maliciousBlocks: Math.floor(Math.random() * 100),
        },
        '2': {
          honestNodes: Math.floor(config.numNodes * (1 - config.maliciousNodeRatio) / config.numShards),
          maliciousNodes: Math.floor(config.numNodes * config.maliciousNodeRatio / config.numShards),
          honestBlocks: Math.floor(Math.random() * 1000),
          maliciousBlocks: Math.floor(Math.random() * 100),
        },
      },
      networkMetrics: {
        averageBlockDelay: {
          '0': Math.random() * 100,
          '1': Math.random() * 100,
          '2': Math.random() * 100,
        },
        averageHeaderDelay: Math.random() * 50,
        averageDownloadDelay: {
          '0': Math.random() * 200,
          '1': Math.random() * 200,
          '2': Math.random() * 200,
        },
      },
    };

    return NextResponse.json(mockResults);
  } catch (error) {
    console.error('Error in simulation:', error);
    return NextResponse.json(
      { error: 'Failed to run simulation' },
      { status: 500 }
    );
  }
}
