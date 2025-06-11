/**
 * Unit tests for getStyleForZoom and the MapStyle component.
 * Framework: Jest with React Testing Library
 */

import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import MapStyle, { getStyleForZoom } from '../MapStyle';
import mapboxgl from 'mapbox-gl';

// Mock mapbox-gl to prevent actual map rendering and side effects
jest.mock('mapbox-gl', () => ({
  Map: jest.fn().mockImplementation(() => ({
    on: jest.fn(),
    off: jest.fn(),
    remove: jest.fn(),
    setStyle: jest.fn(),
    addControl: jest.fn(),
  })),
  NavigationControl: jest.fn(),
  ScaleControl: jest.fn(),
}));

afterEach(() => {
  jest.clearAllMocks();
});

describe('getStyleForZoom', () => {
  test.each([
    { zoom: 0, expected: 'low' },
    { zoom: 5, expected: 'low' },
    { zoom: 6, expected: 'medium' },
    { zoom: 10, expected: 'medium' },
    { zoom: 11, expected: 'high' },
    { zoom: 15, expected: 'high' },
    { zoom: 16, expected: 'ultra' },
    { zoom: 20, expected: 'ultra' },
  ])('returns "%expected" for zoom level $zoom', ({ zoom, expected }) => {
    expect(getStyleForZoom(zoom)).toBe(expected);
  });

  it('throws for invalid zoom values', () => {
    expect(() => getStyleForZoom(-1)).toThrow();
    expect(() => getStyleForZoom(NaN)).toThrow();
  });
});

describe('MapStyle component rendering', () => {
  it('renders map container with light theme at zoom 5', () => {
    render(<MapStyle theme="light" zoom={5} />);
    const container = screen.getByTestId('map-container');
    expect(container).toBeInTheDocument();
    expect(mapboxgl.Map).toHaveBeenCalledTimes(1);
  });

  it('applies dark theme style when theme="dark"', () => {
    render(<MapStyle theme="dark" zoom={8} />);
    const mapInstance = (mapboxgl.Map as jest.Mock).mock.results[0].value;
    expect(mapInstance.setStyle).toHaveBeenCalledWith('dark');
  });

  it('does not crash with missing props', () => {
    // @ts-ignore: intentionally omitting props
    const renderComp = () => render(<MapStyle />);
    expect(renderComp).not.toThrow();
    const container = screen.getByTestId('map-container');
    expect(container).toBeInTheDocument();
  });

  it('catches and logs errors from mapbox-gl constructor', () => {
    const error = new Error('Map init failure');
    (mapboxgl.Map as jest.Mock).mockImplementationOnce(() => { throw error; });
    const spy = jest.spyOn(console, 'error').mockImplementation(() => {});
    expect(() => render(<MapStyle theme="light" zoom={3} />)).not.toThrow();
    expect(spy).toHaveBeenCalledWith(error);
    spy.mockRestore();
  });
});