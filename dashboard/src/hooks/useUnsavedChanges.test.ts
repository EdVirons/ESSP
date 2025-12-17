import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useBeforeUnload } from './useUnsavedChanges';

describe('useBeforeUnload', () => {
  let addEventListenerSpy: ReturnType<typeof vi.spyOn>;
  let removeEventListenerSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    addEventListenerSpy = vi.spyOn(window, 'addEventListener');
    removeEventListenerSpy = vi.spyOn(window, 'removeEventListener');
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should add beforeunload event listener when isDirty is true', () => {
    renderHook(() => useBeforeUnload(true));

    expect(addEventListenerSpy).toHaveBeenCalledWith(
      'beforeunload',
      expect.any(Function)
    );
  });

  it('should add beforeunload event listener even when isDirty is false', () => {
    renderHook(() => useBeforeUnload(false));

    expect(addEventListenerSpy).toHaveBeenCalledWith(
      'beforeunload',
      expect.any(Function)
    );
  });

  it('should remove event listener on unmount', () => {
    const { unmount } = renderHook(() => useBeforeUnload(true));

    unmount();

    expect(removeEventListenerSpy).toHaveBeenCalledWith(
      'beforeunload',
      expect.any(Function)
    );
  });

  it('should prevent default when isDirty and beforeunload is triggered', () => {
    renderHook(() => useBeforeUnload(true));

    // Get the handler function that was registered
    const handler = addEventListenerSpy.mock.calls.find(
      (call) => call[0] === 'beforeunload'
    )?.[1] as EventListener;

    expect(handler).toBeDefined();

    // Create a mock event
    const mockEvent = {
      preventDefault: vi.fn(),
      returnValue: '',
    } as unknown as BeforeUnloadEvent;

    // Trigger the handler
    handler(mockEvent);

    expect(mockEvent.preventDefault).toHaveBeenCalled();
  });

  it('should not prevent default when isDirty is false', () => {
    renderHook(() => useBeforeUnload(false));

    const handler = addEventListenerSpy.mock.calls.find(
      (call) => call[0] === 'beforeunload'
    )?.[1] as EventListener;

    const mockEvent = {
      preventDefault: vi.fn(),
      returnValue: '',
    } as unknown as BeforeUnloadEvent;

    handler(mockEvent);

    expect(mockEvent.preventDefault).not.toHaveBeenCalled();
  });

  it('should update listener when isDirty changes', () => {
    const { rerender } = renderHook(
      ({ isDirty }) => useBeforeUnload(isDirty),
      { initialProps: { isDirty: false } }
    );

    // Initial call
    expect(addEventListenerSpy).toHaveBeenCalledTimes(1);

    // Rerender with different value
    rerender({ isDirty: true });

    // Should have removed old listener and added new one
    expect(removeEventListenerSpy).toHaveBeenCalled();
    expect(addEventListenerSpy).toHaveBeenCalledTimes(2);
  });
});
