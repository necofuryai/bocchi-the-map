// Tests use Jest 29 + @testing-library/react-hooks
import { renderHook, act } from '@testing-library/react-hooks';
import { waitFor } from '@testing-library/react';
import useMaplibre from '../useMaplibre';
import { EventEmitter } from 'events';

jest.mock('maplibre-gl', () => ({
  Map: jest.fn().mockImplementation(() => {
    const emitter = new EventEmitter();
    return {
      on: emitter.on.bind(emitter),
      off: emitter.removeListener.bind(emitter),
      remove: jest.fn(),
      setCenter: jest.fn(),
      setZoom: jest.fn(),
      // expose emitter for manual event firing
      __emitter: emitter,
    };
  }),
}));

describe('useMaplibre', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('initializes the map and sets isLoaded to true', async () => {
    const container = document.createElement('div');
    const ref = { current: container } as React.RefObject<HTMLDivElement>;
    const onLoad = jest.fn();
    const { result } = renderHook(() =>
      useMaplibre({ container: ref, center: [0, 0], zoom: 1, onLoad })
    );
    const mapInstance = (require('maplibre-gl').Map as jest.Mock).mock.results[0].value;
    // simulate load event
    act(() => {
      mapInstance.__emitter.emit('load');
    });
    await waitFor(() => expect(result.current.isLoaded).toBe(true));
    expect(onLoad).toHaveBeenCalledTimes(1);
  });

  it('updates center and zoom when props change', async () => {
    const container = document.createElement('div');
    const ref = { current: container } as any;
    const initial = { container: ref, center: [0, 0], zoom: 1 };
    const updated = { container: ref, center: [10, 10], zoom: 5 };
    const { result, rerender } = renderHook(
      props => useMaplibre(props),
      { initialProps: initial }
    );
    const mapInstance = (require('maplibre-gl').Map as jest.Mock).mock.results[0].value;
    act(() => {
      mapInstance.__emitter.emit('load');
    });
    await waitFor(() => expect(result.current.isLoaded).toBe(true));
    act(() => {
      rerender(updated);
    });
    expect(mapInstance.setCenter).toHaveBeenCalledWith([10, 10]);
    expect(mapInstance.setZoom).toHaveBeenCalledWith(5);
  });

  it('surfaces error when Map constructor throws', async () => {
    const error = new Error('maplibre fail');
    (require('maplibre-gl').Map as jest.Mock).mockImplementationOnce(() => {
      throw error;
    });
    const container = document.createElement('div');
    const ref = { current: container } as any;
    const onError = jest.fn();
    const { result } = renderHook(() =>
      useMaplibre({ container: ref, center: [0, 0], zoom: 1, onError })
    );
    await waitFor(() => expect(result.current.error).toBe(error));
    expect(onError).toHaveBeenCalledWith(error);
  });

  it('does nothing when container ref is null', () => {
    const ref = { current: null } as any;
    const { result } = renderHook(() =>
      useMaplibre({ container: ref, center: [0, 0], zoom: 1 })
    );
    expect(result.current.isLoaded).toBe(false);
    expect(result.current.error).toBeNull();
    expect((require('maplibre-gl').Map as jest.Mock)).not.toHaveBeenCalled();
  });

  it('cleans up map and listeners on unmount', async () => {
    const container = document.createElement('div');
    const ref = { current: container } as any;
    const { result, unmount } = renderHook(() =>
      useMaplibre({ container: ref, center: [0, 0], zoom: 1 })
    );
    const mapInstance = (require('maplibre-gl').Map as jest.Mock).mock.results[0].value;
    act(() => {
      mapInstance.__emitter.emit('load');
    });
    await waitFor(() => expect(result.current.isLoaded).toBe(true));
    unmount();
    expect(mapInstance.remove).toHaveBeenCalled();
    expect(mapInstance.off).toHaveBeenCalledWith('load', expect.any(Function));
    expect(mapInstance.off).toHaveBeenCalledWith('move', expect.any(Function));
    expect(mapInstance.off).toHaveBeenCalledWith('error', expect.any(Function));
  });

  it('creates map instance only once', async () => {
    const container = document.createElement('div');
    const ref = { current: container } as any;
    const { result, rerender } = renderHook(
      props => useMaplibre(props),
      { initialProps: { container: ref, center: [0, 0], zoom: 1 } }
    );
    const mapMock = require('maplibre-gl').Map as jest.Mock;
    act(() => {
      mapMock.mock.results[0].value.__emitter.emit('load');
    });
    await waitFor(() => expect(result.current.isLoaded).toBe(true));
    rerender({ container: ref, center: [0, 0], zoom: 1 });
    expect(mapMock).toHaveBeenCalledTimes(1);
  });
});