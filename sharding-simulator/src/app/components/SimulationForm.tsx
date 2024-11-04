'use client';

import { useState } from 'react';
import {
  Box,
  Text,
  VStack,
  SimpleGrid,
  Heading,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Button,
} from '@chakra-ui/react';
import type { SimulationConfig } from '../types';

interface Props {
  onSubmit: (config: SimulationConfig) => void;
  loading: boolean;
}

export default function SimulationForm({ onSubmit, loading }: Props) {
  const [config, setConfig] = useState<SimulationConfig>({
    numNodes: 10000,
    numOperators: 20,
    numShards: 3,
    simulationTime: 5000,
    timeStep: 1,
    attackStartTime: 20,
    attackEndTime: 60,
    blockProductionInterval: 6,
    transactionsPerBlock: 6500,
    maliciousNodeRatio: 0.1,
    lotteryWinProbability: 0.001,
    maliciousNodeMultiplier: 0,
    blockSize: 115 * 6500,
    blockHeaderSize: 1000,
    erHeaderSize: 1000,
    erBodySize: 33000,
    networkBandwidth: 10,
    minNetworkDelayMean: 50.0,
    maxNetworkDelayMean: 200.0,
    minNetworkDelayStd: 10.0,
    maxNetworkDelayStd: 50.0,
    minGossipFanout: 4,
    maxGossipFanout: 8,
    maxP2PConnections: 2,
    timeOut: 2000,
    numBlocksToDownload: 100
  });

  const [floatInputs, setFloatInputs] = useState({
    lotteryWinProbability: config.lotteryWinProbability.toString(),
    maliciousNodeRatio: config.maliciousNodeRatio.toString()
  });

  const handleNumberChange = (value: number, key: keyof SimulationConfig) => {
    setConfig(prev => ({ ...prev, [key]: value }));
  };

  const renderParameter = (
    label: string,
    value: number,
    unit: string,
    key: keyof SimulationConfig
  ) => {
    const needsFloatHandling = key === 'lotteryWinProbability' || key === 'maliciousNodeRatio';
    
    return (
      <Box>
        <Text color="gray.400" fontSize="sm" mb={2}>{label}</Text>
        {needsFloatHandling ? (
          <NumberInput
            value={floatInputs[key]}
            onChange={(valueString) => {
              setFloatInputs(prev => ({
                ...prev,
                [key]: valueString
              }));
            }}
            bg="gray.800"
            borderRadius="md"
            borderWidth="1px"
            borderColor={['numShards', 'maliciousNodeRatio', 'lotteryWinProbability'].includes(key) ? "blue.400" : "transparent"}
          >
            <NumberInputField
              height="48px"
              color="white"
              fontSize="xl"
              fontWeight="medium"
              border="none"
              _focus={{
                borderColor: "blue.400",
                boxShadow: "none"
              }}
            />
          </NumberInput>
        ) : (
          <NumberInput
            value={value}
            onChange={(_, val) => handleNumberChange(val, key)}
            min={0}
            bg="gray.800"
            borderRadius="md"
            borderWidth="1px"
            borderColor={['numShards', 'maliciousNodeRatio', 'lotteryWinProbability'].includes(key) ? "blue.400" : "transparent"}
          >
            <NumberInputField
              height="48px"
              color="white"
              fontSize="xl"
              fontWeight="medium"
              border="none"
              _focus={{
                borderColor: "blue.400",
                boxShadow: "none"
              }}
            />
            <NumberInputStepper>
              <NumberIncrementStepper borderColor="gray.700" color="gray.400" />
              <NumberDecrementStepper borderColor="gray.700" color="gray.400" />
            </NumberInputStepper>
          </NumberInput>
        )}
        <Text color="gray.500" fontSize="xs" mt={1}>
          {`${value} ${unit}`}
        </Text>
      </Box>
    );
  };

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      
      const finalConfig = {
        ...config,
        lotteryWinProbability: parseFloat(floatInputs.lotteryWinProbability) || 0,
        maliciousNodeRatio: parseFloat(floatInputs.maliciousNodeRatio) || 0
      };
      
      onSubmit(finalConfig);
    }}>
      <VStack spacing={8} align="stretch">
        {/* Node Configuration */}
        <Box bg="gray.900" p={6} borderRadius="lg">
          <Heading size="md" color="white" mb={6}>Node Configuration</Heading>
          <SimpleGrid columns={{ base: 2, md: 4 }} spacing={6}>
            {renderParameter('Network Size', config.numNodes, 'nodes', 'numNodes')}
            {renderParameter('Operators', config.numOperators, 'operators', 'numOperators')}
            {renderParameter('Shards', config.numShards, 'shards', 'numShards')}
            {renderParameter('Time Step', config.timeStep, 'units', 'timeStep')}
            {renderParameter('Simulation Time', config.simulationTime, 'time units', 'simulationTime')}
            {renderParameter('Block Interval', config.blockProductionInterval, 'blocks', 'blockProductionInterval')}
            {renderParameter('Tx Per Block', config.transactionsPerBlock, 'tx', 'transactionsPerBlock')}
            {renderParameter('Lottery Win Prob', config.lotteryWinProbability, '%', 'lotteryWinProbability')}
          </SimpleGrid>
        </Box>

        {/* Block Parameters */}
        <Box bg="gray.900" p={6} borderRadius="lg">
          <Heading size="md" color="white" mb={6}>Block Parameters</Heading>
          <SimpleGrid columns={{ base: 2, md: 4 }} spacing={6}>
            {renderParameter('Block Header Size', config.blockHeaderSize, 'bytes', 'blockHeaderSize')}
            {/* {renderParameter('ER Header Size', config.erHeaderSize, 'bytes', 'erHeaderSize')}
            {renderParameter('ER Body Size', config.erBodySize, 'bytes', 'erBodySize')} */}
            {renderParameter('Blocks to Download', config.numBlocksToDownload, 'blocks', 'numBlocksToDownload')}
          </SimpleGrid>
        </Box>

        {/* Network Parameters */}
        <Box bg="gray.900" p={6} borderRadius="lg">
          <Heading size="md" color="white" mb={6}>Network Parameters</Heading>
          <SimpleGrid columns={{ base: 2, md: 4 }} spacing={6}>
            {renderParameter('Network Bandwidth', config.networkBandwidth, 'Mbps', 'networkBandwidth')}
            {renderParameter('Min Network Delay Mean', config.minNetworkDelayMean, 'ms', 'minNetworkDelayMean')}
            {renderParameter('Max Network Delay Mean', config.maxNetworkDelayMean, 'ms', 'maxNetworkDelayMean')}
            {renderParameter('Min Network Delay Std', config.minNetworkDelayStd, 'ms', 'minNetworkDelayStd')}
            {renderParameter('Max Network Delay Std', config.maxNetworkDelayStd, 'ms', 'maxNetworkDelayStd')}
            {renderParameter('Min Gossip Fanout', config.minGossipFanout, 'peers', 'minGossipFanout')}
            {renderParameter('Max Gossip Fanout', config.maxGossipFanout, 'peers', 'maxGossipFanout')}
            {renderParameter('Max P2P Connections', config.maxP2PConnections, 'connections', 'maxP2PConnections')}
            {renderParameter('Timeout', config.timeOut, 'ms', 'timeOut')}
          </SimpleGrid>
        </Box>

        {/* Attack Configuration */}
        <Box bg="gray.900" p={6} borderRadius="lg">
          <Heading size="md" color="white" mb={6}>Attack Configuration</Heading>
          <SimpleGrid columns={{ base: 2, md: 4 }} spacing={6}>
            {renderParameter('Malicious Node Ratio', (config.maliciousNodeRatio) * 100, '%', 'maliciousNodeRatio')}
            {/* {renderParameter('Malicious Node Multiplier', config.maliciousNodeMultiplier, 'x', 'maliciousNodeMultiplier')} */}
            {/* {renderParameter('Attack Start Time', config.attackStartTime, 'time units', 'attackStartTime')}
            {renderParameter('Attack End Time', config.attackEndTime, 'time units', 'attackEndTime')} */}
          </SimpleGrid>
        </Box>

        <Button
          type="submit"
          isLoading={loading}
          colorScheme="blue"
          size="lg"
          height="54px"
          fontSize="lg"
          width="full"
          bg="blue.600"
          _hover={{ bg: 'blue.500' }}
          _active={{ bg: 'blue.700' }}
        >
          {loading ? 'Running Simulation...' : 'Run Simulation'}
        </Button>
      </VStack>
    </form>
  );
}