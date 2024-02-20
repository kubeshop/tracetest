import _get from 'lodash/get';
import type React from 'react';

type TNil = null | undefined;

enum EUpdateTypes {
  DragEnd = 'DragEnd',
  DragMove = 'DragMove',
  DragStart = 'DragStart',
}

type TDraggableBounds = {
  clientXLeft: number;
  maxValue?: number;
  minValue?: number;
  width: number;
};

type TDraggableUpdate = {
  event: React.MouseEvent<HTMLDivElement | SVGSVGElement, MouseEvent> | MouseEvent;
  type: EUpdateTypes;
  value: number;
  x: number;
  resetBounds(): void;
};

type TDraggableManagerOptions = {
  getBounds: () => TDraggableBounds;
  onDragStart?: (update: TDraggableUpdate) => void;
  onDragMove?: (update: TDraggableUpdate) => void;
  onDragEnd?: (update: TDraggableUpdate) => void;
};

interface IDraggableManager {
  init(): (event: MouseEvent) => void;
  resetBounds(): void;
  dispose(): void;
}

const LEFT_MOUSE_BUTTON = 0;

function DraggableManager({getBounds, ...rest}: TDraggableManagerOptions): IDraggableManager {
  let _bounds: TDraggableBounds | TNil;
  let _isDragging: boolean;
  let _onDragStart: ((update: TDraggableUpdate) => void) | TNil;
  let _onDragMove: ((update: TDraggableUpdate) => void) | TNil;
  let _onDragEnd: ((update: TDraggableUpdate) => void) | TNil;

  const _resetBounds = () => {
    _bounds = undefined;
  };

  const _getBounds = () => {
    if (!_bounds) {
      _bounds = getBounds();
    }
    return _bounds;
  };

  const _getPosition = (clientX: number) => {
    const {clientXLeft, maxValue, minValue, width} = _getBounds();
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

  const _stopDragging = () => {
    window.removeEventListener('mousemove', _handleDragEvent);
    window.removeEventListener('mouseup', _handleDragEvent);
    const style = _get(document, 'body.style');
    if (style) {
      style.removeProperty('userSelect');
    }
    _isDragging = false;
  };

  const _handleDragEvent = (event: MouseEvent | React.MouseEvent<HTMLDivElement | SVGSVGElement, MouseEvent>) => {
    const {button, clientX, type: eventType} = event;
    let type: EUpdateTypes | null = null;
    let handler: ((update: TDraggableUpdate) => void) | TNil;

    if (eventType === 'mousedown') {
      if (_isDragging || button !== LEFT_MOUSE_BUTTON) {
        return;
      }

      window.addEventListener('mousemove', _handleDragEvent);
      window.addEventListener('mouseup', _handleDragEvent);

      const style = _get(document, 'body.style');
      if (style) {
        style.userSelect = 'none';
      }
      _isDragging = true;

      type = EUpdateTypes.DragStart;
      handler = _onDragStart;
    } else if (eventType === 'mousemove') {
      if (!_isDragging) {
        return;
      }

      type = EUpdateTypes.DragMove;
      handler = _onDragMove;
    } else if (eventType === 'mouseup') {
      if (!_isDragging) {
        return;
      }

      _stopDragging();
      type = EUpdateTypes.DragEnd;
      handler = _onDragEnd;
    } else {
      throw new Error(`invalid event type: ${eventType}`);
    }

    if (!handler) {
      return;
    }

    const {value, x} = _getPosition(clientX);

    handler({
      event,
      type,
      value,
      x,
      resetBounds: _resetBounds,
    });
  };

  return {
    init() {
      _isDragging = false;
      _bounds = undefined;
      window.addEventListener('resize', this.resetBounds);
      _onDragStart = rest.onDragStart;
      _onDragMove = rest.onDragMove;
      _onDragEnd = rest.onDragEnd;

      return _handleDragEvent;
    },

    resetBounds: _resetBounds,

    dispose() {
      if (_isDragging) {
        _stopDragging();
      }
      window.removeEventListener('resize', this.resetBounds);
      _bounds = undefined;
      _onDragStart = undefined;
      _onDragMove = undefined;
      _onDragEnd = undefined;
    },
  };
}

export default DraggableManager;
