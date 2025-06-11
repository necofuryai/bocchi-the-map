afterEach(() => {
  cleanup();
});

describe('POIFeatures - Additional Unit Tests', () => {
  describe('Happy Path Rendering', () => {
    it('renders POI features on map with valid GeoJSON', async () => {
      const validGeoJSON = {
        type: 'FeatureCollection',
        features: [
          {
            type: 'Feature',
            geometry: { type: 'Point', coordinates: [100, 0] },
            properties: { id: 'poi1', name: 'POI 1' },
          },
        ],
      };
      render(<POIFeatures data={validGeoJSON} />);
      const poiLabel = await screen.findByText('POI 1');
      expect(poiLabel).toBeInTheDocument();
    });
  });

  describe('Edge Cases', () => {
    it('renders no features when provided an empty feature list', () => {
      const emptyGeoJSON = { type: 'FeatureCollection', features: [] };
      render(<POIFeatures data={emptyGeoJSON} />);
      expect(screen.queryByTestId('poi-marker')).toBeNull();
    });

    it('skips rendering features with null geometry', () => {
      const nullGeoJSON = {
        type: 'FeatureCollection',
        features: [
          { type: 'Feature', geometry: null, properties: { id: 'poiNull', name: 'Null POI' } },
        ],
      };
      render(<POIFeatures data={nullGeoJSON} />);
      expect(screen.queryByText('Null POI')).toBeNull();
    });

    it('clusters overlapping features into a cluster marker', () => {
      const overlappingGeoJSON = {
        type: 'FeatureCollection',
        features: [
          {
            type: 'Feature',
            geometry: { type: 'Point', coordinates: [0, 0] },
            properties: { id: 'poi1', name: 'POI 1' },
          },
          {
            type: 'Feature',
            geometry: { type: 'Point', coordinates: [0, 0] },
            properties: { id: 'poi2', name: 'POI 2' },
          },
        ],
      };
      render(<POIFeatures data={overlappingGeoJSON} />);
      const cluster = screen.getByTestId('poi-cluster');
      expect(cluster).toBeInTheDocument();
    });
  });

  describe('Failure Conditions', () => {
    beforeEach(() => {
      jest.spyOn(console, 'error').mockImplementation(() => {});
    });

    afterEach(() => {
      (console.error as jest.Mock).mockRestore();
    });

    it('logs an error for invalid coordinates without throwing', () => {
      const invalidGeoJSON = {
        type: 'FeatureCollection',
        features: [
          {
            type: 'Feature',
            geometry: { type: 'Point', coordinates: ['invalid', 'coords'] as any },
            properties: { id: 'badPoi', name: 'Bad POI' },
          },
        ],
      };
      render(<POIFeatures data={invalidGeoJSON} />);
      expect(console.error).toHaveBeenCalled();
    });

    it('displays an error message when data fetch service fails', async () => {
      const mockFetchPOIs = jest.fn().mockRejectedValue(new Error('Network Error'));
      jest.doMock('../../services/poiService', () => ({ fetchPOIs: mockFetchPOIs }));
      // Re-require component after mocking service
      const { default: POIFeaturesWithError } = require('../poi-features');
      render(<POIFeaturesWithError />);
      const errorMsg = await screen.findByText(/failed to load POIs/i);
      expect(errorMsg).toBeInTheDocument();
    });
  });
});