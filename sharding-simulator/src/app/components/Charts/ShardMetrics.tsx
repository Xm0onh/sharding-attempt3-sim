'use client';

import { BarChart } from '@tremor/react';
import { SimulationResults } from '../../types';

interface Props {
  blockProduction: SimulationResults['block_production'];
}

export default function ShardMetrics({ blockProduction }: Props) {
  if (!blockProduction || Object.keys(blockProduction).length === 0) {
    return <div>No shard metrics available</div>;
  }

  const data = Object.entries(blockProduction).map(([shardId, stats]) => ({
    shard: `Shard ${shardId}`,
    'Honest Blocks': stats.honest_blocks,
    'Malicious Blocks': stats.malicious_blocks,
  }));

  return (
    <BarChart
      data={data}
      index="shard"
      categories={['Honest Blocks', 'Malicious Blocks']}
      colors={['emerald', 'red']}
      yAxisWidth={56}
      className="h-80 mt-4"
      showLegend={true}
      valueFormatter={(number: number) => Intl.NumberFormat('us').format(number).toString()}
    />
  );
}
