'use client';

import { useState, useEffect } from 'react';
import SimulationForm from './components/SimulationForm';
import SimulationResults from './components/SimulationResults';
import type { SimulationConfig, SimulationResults as Results } from './types';
import { Card } from '@tremor/react';

export default function Home() {
  const [results, setResults] = useState<Results | null>(null);
  const [loading, setLoading] = useState(false);

  const runSimulation = async (config?: SimulationConfig) => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8080/simulate', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        }
      });
      const data = await response.json();
      console.log('Received data:', data);
      setResults(data);
    } catch (error) {
      console.error('Error fetching simulation results:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    runSimulation();
  }, []);

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
