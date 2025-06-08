import { useEffect } from 'react';
import maplibregl from 'maplibre-gl';
import { Protocol } from 'pmtiles';

export const usePmtiles = () => {
  useEffect(() => {
    try {
      const protocol = new Protocol();
      maplibregl.addProtocol('pmtiles', protocol.tile);
      
      return () => {
        maplibregl.removeProtocol('pmtiles');
      };
    } catch (error) {
      console.error('Failed to initialize PMTiles protocol:', error);
    }
  }, []);
};