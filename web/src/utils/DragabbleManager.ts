import _get from 'lodash/get';
import type React from 'react';

type TNil = null | undefined;

type TDraggableBounds = {
  clientXLeft: number;
  maxValue?: number;
  minValue?: number;
  width: number;
};

type TDraggableUpdate = {
  event: React.MouseEvent<HTMLDivElement | SVGSVGElement, MouseEvent> | MouseEvent;
  resetBounds(): void;
  value: number;
  x: number;
};

type TDraggableManagerOptions = {
  onGetBounds: () => TDraggableBounds;
  onDragEnd?: (update: TDraggableUpdate) => void;
  onDragMove?: (update: TDraggableUpdate) => void;
  onDragStart?: (update: TDraggableUpdate) => void;
};

interface IDraggableManager {
  initEventHandler(): (event: MouseEvent) => void;
  cleanup(): void;
  resetBounds(): void;
}

const LEFT_MOUSE_BUTTON = 0;

function DraggableManager({onGetBounds, ...rest}: TDraggableManagerOptions): IDraggableManager {
  let bounds: TDraggableBounds | TNil;
  let isDragging: boolean;
  let onDragStart: ((update: TDraggableUpdate) => void) | TNil;
  let onDragMove: ((update: TDraggableUpdate) => void) | TNil;
  let onDragEnd: ((update: TDraggableUpdate) => void) | TNil;

  const resetBounds = () => {
    bounds = undefined;
  };

  const getBounds = () => {
    if (!bounds) {
      bounds = onGetBounds();
    }
    return bounds;
  };

  const getPosition = (clientX: number) => {
    const {clientXLeft, maxValue, minValue, width} = getBounds();
    let x = clientX - clientXLeft;
    let value = x / width;
    if (minValue != null && value < minValue) {
      value = minValue;
      x = minValue * width;
    } else if (maxValue != null && value > maxValue) {
      value = maxValue;
      x = maxValue * width;
    }
    return {value, x};
  };

  const stopDragging = () => {
    window.removeEventListener('mousemove', handleEvent);
    window.removeEventListener('mouseup', handleEvent);
    const style = _get(document, 'body.style');
    if (style) {
      style.removeProperty('userSelect');
    }
    isDragging = false;
  };

  const handleEvent = (event: MouseEvent | React.MouseEvent<HTMLDivElement | SVGSVGElement, MouseEvent>) => {
    const {button, clientX, type: eventType} = event;
    let handler: ((update: TDraggableUpdate) => void) | TNil;

    if (eventType === 'mousedown') {
      if (isDragging || button !== LEFT_MOUSE_BUTTON) {
        return;
      }

      window.addEventListener('mousemove', handleEvent);
      window.addEventListener('mouseup', handleEvent);

      const style = _get(document, 'body.style');
      if (style) {
        style.userSelect = 'none';
      }
      isDragging = true;

      handler = onDragStart;
    } else if (eventType === 'mousemove') {
      if (!isDragging) {
        return;
      }

      handler = onDragMove;
    } else if (eventType === 'mouseup') {
      if (!isDragging) {
        return;
      }

      stopDragging();
      handler = onDragEnd;
    } else {
      throw new Error(`invalid event type: ${eventType}`);
    }

    if (!handler) {
      return;
    }

    const {value, x} = getPosition(clientX);

    handler({
      event,
      value,
      x,
      resetBounds,
    });
  };

  return {
    initEventHandler() {
      isDragging = false;
      bounds = undefined;
      window.addEventListener('resize', this.resetBounds);
      onDragStart = rest.onDragStart;
      onDragMove = rest.onDragMove;
      onDragEnd = rest.onDragEnd;

      return handleEvent;
    },

    cleanup() {
      if (isDragging) {
        stopDragging();
      }
      window.removeEventListener('resize', this.resetBounds);
      bounds = undefined;
      onDragStart = undefined;
      onDragMove = undefined;
      onDragEnd = undefined;
    },

    resetBounds,
  };
}

export default DraggableManager;
