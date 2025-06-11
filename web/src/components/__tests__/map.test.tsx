// TODO: Add test for initial loading indicator while map script is loading
// TODO: Add test for map initialization failure (calls onError)
// TODO: Add test for prop changes (center/zoom updates)
// TODO: Add test for cleaning up event listeners on unmount
// TODO: Add accessibility test using jest-axe

// Mock global Google Maps API for all tests
beforeAll(() => {
  (window as any).google = {
    maps: {
      Map: jest.fn().mockImplementation(function () {
        this.setCenter = jest.fn();
        this.setZoom = jest.fn();
        this.addListener = jest.fn().mockReturnValue('resizeListener');
      }),
      event: {
        removeListener: jest.fn(),
      },
    },
  };
});

// Ensure clean slate after each test
afterEach(() => {
  jest.clearAllMocks();
  cleanup();
});

describe('Map component', () => {
  it('renders loader while map script is loading', async () => {
    // Simulate script not yet loaded
    (window as any).google = undefined;
    render(<Map center={{ lat: 0, lng: 0 }} zoom={5} onError={jest.fn()} />);
    expect(screen.getByRole('progressbar')).toBeInTheDocument();

    // Simulate script load
    act(() => {
      (window as any).google = {
        maps: (window as any).google?.maps,
      };
    });
    await waitFor(() => {
      expect(screen.queryByRole('progressbar')).not.toBeInTheDocument();
    });
  });

  it('calls onError when map fails to initialize', async () => {
    const initError = new Error('init failure');
    // Make the first Map constructor throw
    (window as any).google.maps.Map.mockImplementationOnce(() => {
      throw initError;
    });
    const onError = jest.fn();
    render(<Map center={{ lat: 0, lng: 0 }} zoom={5} onError={onError} />);
    await waitFor(() => {
      expect(onError).toHaveBeenCalledWith(initError);
    });
  });

  it('updates view when center/zoom props change', async () => {
    const onError = jest.fn();
    const { rerender } = render(
      <Map center={{ lat: 0, lng: 0 }} zoom={5} onError={onError} />
    );
    const mapInstance = (window as any).google.maps.Map.mock.instances[0];

    rerender(<Map center={{ lat: 1, lng: 2 }} zoom={10} onError={onError} />);
    await waitFor(() => {
      expect(mapInstance.setCenter).toHaveBeenCalledWith({ lat: 1, lng: 2 });
      expect(mapInstance.setZoom).toHaveBeenCalledWith(10);
    });
  });

  it('removes resize listener on unmount', () => {
    const onError = jest.fn();
    const { unmount } = render(
      <Map center={{ lat: 0, lng: 0 }} zoom={5} onError={onError} />
    );
    const removeListener = (window as any).google.maps.event.removeListener;
    unmount();
    expect(removeListener).toHaveBeenCalledWith('resizeListener');
  });

  it('should have no accessibility violations', async () => {
    // Ensure `import { axe } from 'jest-axe'` is at the top of this file
    const { container } = render(
      <Map center={{ lat: 0, lng: 0 }} zoom={5} onError={jest.fn()} />
    );
    const results = await axe(container);
    expect(results).toHaveNoViolations();
  });
});