'use client';

import React, { useState, useEffect } from 'react';

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

const TrainList = () => {
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
              // Try to reconnect in 3 seconds
              reconnectTimeout = setTimeout(connect, 3000);
          };

          ws.onerror = (error) => {
              console.error('WebSocket error:', error);
              setConnectionStatus('error');
          };
      };

      connect();

      // Cleanup on component unmount
      return () => {
          if (ws) {
              ws.close();
          }
          if (reconnectTimeout) {
              clearTimeout(reconnectTimeout);
          }
      };
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="bg-white rounded-lg shadow">
                <div className="px-4 py-5 sm:p-6">
                    <div className="flex justify-between items-center mb-6">
                        <h1 className="text-3xl font-bold text-gray-900">Victoria Line Trains</h1>
                        <div className={`px-3 py-1 rounded-full text-sm ${
                            connectionStatus === 'connected' ? 'bg-green-100 text-green-800' :
                            connectionStatus === 'reconnecting' ? 'bg-yellow-100 text-yellow-800' :
                            'bg-red-100 text-red-800'
                        }`}>
                            {connectionStatus}
                        </div>
                    </div>

                        <div className="overflow-x-auto">
                            <table className="min-w-full divide-y divide-gray-200">
                                <thead>
                                    <tr>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                            Train ID
                                        </th>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                            Location
                                        </th>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                            Direction
                                        </th>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                            Next Station In
                                        </th>
                                    </tr>
                                </thead>
                                <tbody className="bg-white divide-y divide-gray-200">
                                    {Object.entries(trains).map(([id, train]) => (
                                        <tr key={id} className="hover:bg-gray-50">
                                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                                {id}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                {train.Location.IsBetween
                                                    ? `Between ${train.Location.PrevStationID} and ${train.Location.StationID}`
                                                    : `${train.Location.State === 'APPROACHING' ? 'Approaching' : 'At'} ${train.Location.StationID}`
                                                }
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                {train.Direction}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                {Math.floor(train.TimeToNext / 60)}m {train.TimeToNext % 60}s
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>

                        <div className="mt-4 text-sm text-gray-500">
                            Active Trains: {Object.keys(trains).length}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TrainList;