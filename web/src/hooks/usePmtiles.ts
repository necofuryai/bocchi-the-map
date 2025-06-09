import { useEffect } from 'react';
import maplibregl from 'maplibre-gl';
import { Protocol } from 'pmtiles';

// Protocol reference counting to handle multiple map instances
let protocolRefCount = 0;
let protocolInstance: Protocol | null = null;

export const usePmtiles = () => {
  useEffect(() => {
    try {
      // Register protocol only if it's the first instance
      if (protocolRefCount === 0) {
        protocolInstance = new Protocol();
        maplibregl.addProtocol('pmtiles', protocolInstance.tile.bind(protocolInstance));
      }
      protocolRefCount++;
      
      return () => {
        // Only remove protocol when the last instance unmounts
        protocolRefCount--;
        if (protocolRefCount === 0 && protocolInstance) {
          maplibregl.removeProtocol('pmtiles');
          protocolInstance = null;
        }
      };
    } catch (error) {
      console.error('Failed to initialize PMTiles protocol:', error);
    }
  }, []);
};