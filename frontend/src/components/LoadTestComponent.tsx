import React, { useState } from 'react';
import {LoadTest} from "../../wailsjs/go/main/App";
import './LoadTestComponent.css'; // Import the CSS file

interface Metrics {
    MaxLatency: number;
    MinLatency: number;
    TotalLatency: number;
    TotalRequests: number;
    TotalResponses: number;
}

const LoadTestComponent: React.FC = () => {
    const [loadTestConfig, setLoadTestConfig] = useState({
        URL: '',
        RequestsPerSec: 0,
        MaxRequests: 0,
        RequestType: '',
        BearerToken: '',
    });
    const [loadTestStatus, setLoadTestStatus] = useState<'idle' | 'running' | 'done'>('idle');
    const [metrics, setMetrics] = useState<Metrics | null>(null);

    const handleInputChange = (event: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        let value: string | number = event.target.value;
    
        if (event.target.name === 'Url' && !value.startsWith('http://') && !value.startsWith('https://')) {
            value = 'https://' + value;
        }
    
        if (event.target.name === 'RequestsPerSec' || event.target.name === 'MaxRequests') {
            value = parseInt(value, 10);
            if (isNaN(value)) {
                return;
            }
        }
    
        setLoadTestConfig({
            ...loadTestConfig,
            [event.target.name]: value,
        });
    };

    const handleLoadTest = async () => {
        setLoadTestStatus('running');
        try {
            const metrics = await LoadTest(loadTestConfig);
            setLoadTestStatus('done');
            setMetrics(metrics);
            console.log(metrics);
        } catch (err) {
            console.error('Error performing load test:', err);
            setLoadTestStatus('idle');
        }
    };

    return (
        <div className="load-test-form">
            <input type="text" name="URL" placeholder="URL" onChange={handleInputChange} />
            <input type="number" name="RequestsPerSec" placeholder="Requests per second" onChange={handleInputChange} />
            <input type="number" name="MaxRequests" placeholder="Max requests" onChange={handleInputChange} />
            <select name="RequestType" onChange={handleInputChange}>
                <option value="">Select request type</option>
                <option value="GET">GET</option>
                <option value="POST">POST</option>
                <option value="PUT">PUT</option>
                <option value="DELETE">DELETE</option>
                <option value="PATCH">PATCH</option>
            </select>
            <input type="text" name="BearerToken" placeholder="Bearer token" onChange={handleInputChange} />
            <button onClick={handleLoadTest} disabled={loadTestStatus === 'running'}>
                {loadTestStatus === 'running' ? 'Running...' : 'Start Load Test'}
            </button>
            {loadTestStatus === 'done' && metrics && (
                <div className="metrics">
                    <h2>Metrics</h2>
                    <p>Max Latency: {metrics.MaxLatency}ms</p>
                    <p>Min Latency: {metrics.MinLatency}ms</p>
                    <p>Total Latency: {metrics.TotalLatency}ms</p>
                    <p>Total Requests: {metrics.TotalRequests}</p>
                    <p>Total Responses: {metrics.TotalResponses}</p>
                </div>
            )}
        </div>
    );
};

export default LoadTestComponent;