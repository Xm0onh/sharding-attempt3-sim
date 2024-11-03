'use client';

import { BarChart } from '@tremor/react';
import { SimulationResults } from '../../types';

interface Props {
  blockProduction: SimulationResults['block_production'];
  networkMetrics: SimulationResults['network_metrics'];
}

export default function TpsChart({ blockProduction, networkMetrics }: Props) {
  if (!blockProduction || Object.keys(blockProduction).length === 0) {
    return <div>No performance data available</div>;
  }

  const data = Object.entries(blockProduction).map(([shardId, stats]) => {
    const broadcastDelay = networkMetrics.block_broadcast_delays_ms[shardId] || 0;
    const downloadDelay = networkMetrics.block_download_delays_ms[shardId] || 0;
    return {
      shard: `Shard ${shardId}`,
      'Block Production': stats.total_blocks,
      'Broadcast Delay': Math.round(broadcastDelay),
      'Download Delay': Math.round(downloadDelay / 100),
    };
  });

  return (
    <BarChart
      data={data}
      index="shard"
      categories={['Block Production', 'Broadcast Delay', 'Download Delay']}
      colors={['blue', 'yellow', 'red']}
      yAxisWidth={56}
      className="h-80 mt-4"
      showLegend={true}
      valueFormatter={(number: number) => Intl.NumberFormat('us').format(number).toString()}
    />
  );
}
