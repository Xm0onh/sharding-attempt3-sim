'use client';

import { Box, SimpleGrid, Text, Heading, Stat, StatLabel, StatNumber, StatGroup } from '@chakra-ui/react';
import type { SimulationResults } from '../types';

interface Props {
  results: SimulationResults;
}

export default function SimulationResults({ results }: Props) {
  return (
    <Box>
      {/* Main Performance Metrics */}
      <Box bg="gray.900" p={6} borderRadius="lg" mb={6}>
        <SimpleGrid columns={{ base: 1, md: 3 }} spacing={6}>
          <Stat>
            <StatLabel color="blue.400">TPS</StatLabel>
            <StatNumber color="white" fontSize="3xl">
              {results?.performance?.transactions_per_second.toFixed(2)}
            </StatNumber>
          </Stat>
          <Stat>
            <StatLabel color="green.400">Transaction Size</StatLabel>
            <StatNumber color="white" fontSize="3xl">
              {results.transaction_size_bytes}
              <Text as="span" fontSize="sm" color="gray.400" ml={2}>bytes</Text>
            </StatNumber>
          </Stat>
          <Stat>
            <StatLabel color="purple.400">Block Size</StatLabel>
            <StatNumber color="white" fontSize="3xl">
              {results.block_size_kb}
              <Text as="span" fontSize="sm" color="gray.400" ml={2}>KB</Text>
            </StatNumber>
          </Stat>
        </SimpleGrid>
      </Box>

      {/* Block Production Stats */}
      <Box bg="gray.900" p={6} borderRadius="lg" mb={6}>
        <Heading size="md" color="gray.300" mb={4}>Block Production by Shard</Heading>
        <SimpleGrid columns={{ base: 1, md: 3 }} spacing={4}>
          {Object.entries(results.block_production).map(([shardId, stats]) => (
            <Box key={shardId} p={4} borderRadius="lg" border="1px" borderColor="gray.700">
              <Text color="gray.400" fontSize="sm" fontWeight="medium">Shard {shardId}</Text>
              <Box mt={3}>
                <SimpleGrid columns={2} spacing={2}>
                  <Text color="green.400">Honest Blocks</Text>
                  <Text color="white" textAlign="right">{stats.honest_blocks}</Text>
                  <Text color="red.400">Malicious Blocks</Text>
                  <Text color="white" textAlign="right">{stats.malicious_blocks}</Text>
                  <Text color="gray.400" pt={2} borderTop="1px" borderColor="gray.700">Total Blocks</Text>
                  <Text color="white" textAlign="right" pt={2} borderTop="1px" borderColor="gray.700">{stats.total_blocks}</Text>
                </SimpleGrid>
              </Box>
            </Box>
          ))}
        </SimpleGrid>
      </Box>

      {/* Network Metrics */}
      <Box bg="gray.900" p={6} borderRadius="lg">
        <Heading size="md" color="gray.300" mb={4}>Network Metrics</Heading>
        <Box>
          <Box p={4} bg="gray.800" borderRadius="lg" mb={6}>
            <Text color="gray.400" fontSize="sm">Block Header Delay</Text>
            <Text color="white" fontSize="2xl" mt={2}>
              {results.network_metrics.block_header_delay_ms.toFixed(2)}
              <Text as="span" fontSize="sm" color="gray.400" ml={2}>ms</Text>
            </Text>
          </Box>

          <Box mb={6}>
            <Text color="gray.400" fontSize="sm" mb={3}>Broadcast Delays</Text>
            <SimpleGrid columns={{ base: 1, md: 3 }} spacing={4}>
              {Object.entries(results.network_metrics.block_broadcast_delays_ms).map(([shardId, delay]) => (
                <Box key={`broadcast-${shardId}`} p={4} borderRadius="lg" border="1px" borderColor="gray.700">
                  <Text color="gray.400" fontSize="sm">Shard {shardId}</Text>
                  <Text color="white" fontSize="xl" mt={1}>
                    {delay.toFixed(2)}
                    <Text as="span" fontSize="sm" color="gray.400" ml={2}>ms</Text>
                  </Text>
                </Box>
              ))}
            </SimpleGrid>
          </Box>

          <Box>
            <Text color="gray.400" fontSize="sm" mb={3}>Download Delays</Text>
            <SimpleGrid columns={{ base: 1, md: 3 }} spacing={4}>
              {Object.entries(results.network_metrics.block_download_delays_ms).map(([shardId, delay]) => (
                <Box key={`download-${shardId}`} p={4} borderRadius="lg" border="1px" borderColor="gray.700">
                  <Text color="gray.400" fontSize="sm">Shard {shardId}</Text>
                  <Text color="white" fontSize="xl" mt={1}>
                    {delay.toFixed(2)}
                    <Text as="span" fontSize="sm" color="gray.400" ml={2}>ms</Text>
                  </Text>
                </Box>
              ))}
            </SimpleGrid>
          </Box>
        </Box>
      </Box>
    </Box>
  );
}
