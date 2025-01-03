// hooks/useTrains.ts
'use client';

import { useState, useEffect } from 'react';

interface Location {
    StationID: string;
    IsBetween: boolean;
    PrevStationID: string;
    State: 'AT_STATION' | 'BETWEEN' | 'APPROACHING';
}

interface TrainInfo {
    Location: Location;
    Direction: string;
    TimeToNext: number;
}

interface TrainMap {
    [key: string]: TrainInfo;
}

export function useTrains() {
    const [trains, setTrains] = useState<TrainMap>({});
    const [connectionStatus, setConnectionStatus] = useState('connecting');

    useEffect(() => {
        let ws: WebSocket | null = null;
        let reconnectTimeout: NodeJS.Timeout;

        const connect = () => {
            ws = new WebSocket('ws://localhost:8080/ws');

            ws.onopen = () => {
                console.log('Connected to WebSocket');
                setConnectionStatus('connected');
            };

            ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    setTrains(data);
                } catch (error) {
                    console.error('Error parsing websocket data:', error);
                }
            };

            ws.onclose = () => {
                console.log('WebSocket closed, attempting to reconnect...');
                setConnectionStatus('reconnecting');
                reconnectTimeout = setTimeout(connect, 3000);
            };

            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                setConnectionStatus('error');
            };
        };

        connect();

        return () => {
            if (ws) {
                ws.close();
            }
            if (reconnectTimeout) {
                clearTimeout(reconnectTimeout);
            }
        };
    }, []);

    return { trains, connectionStatus };
}