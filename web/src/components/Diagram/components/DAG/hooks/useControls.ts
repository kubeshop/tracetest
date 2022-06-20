import {useCallback, useEffect} from 'react';
import {useStore, useZoomPanHelper} from 'react-flow-renderer';
import {useSpan} from 'providers/Span/Span.provider';

interface IProps {
  onSelectSpan(spanId: string): void;
}

const useControls = ({onSelectSpan}: IProps) => {
  const {setCenter, fitView} = useZoomPanHelper();
  const {focusedSpan, onSetFocusedSpan, affectedSpans} = useSpan();
  const {getState} = useStore();
  const indexOfFocused = affectedSpans.findIndex(spanId => spanId === focusedSpan) + 1;

  const getNodePosition = useCallback(
    (nodeId: string) => {
      const {nodes} = getState();
      const {position} = nodes.find(node => node.id === nodeId) || {};

      return position || {x: 0, y: 0};
    },
    [getState]
  );

  useEffect(() => {
    if (affectedSpans.length && !focusedSpan) {
      const [spanId] = affectedSpans;

      onSelectSpan(spanId);
      onSetFocusedSpan(spanId);
    }
  }, [affectedSpans, focusedSpan, onSelectSpan, onSetFocusedSpan]);

  useEffect(() => {
    if (focusedSpan) {
      const {x, y} = getNodePosition(focusedSpan);

      setCenter(x, y);
    } else fitView();
  }, [fitView, focusedSpan, getNodePosition, getState, setCenter]);

  const handleNextSpan = useCallback(() => {
    const nextSpan = affectedSpans[affectedSpans.indexOf(focusedSpan) + 1] || affectedSpans[0];
    onSelectSpan(nextSpan);
    onSetFocusedSpan(nextSpan);
  }, [affectedSpans, focusedSpan, onSelectSpan, onSetFocusedSpan]);

  const handlePrevSpan = useCallback(() => {
    const preSpan = affectedSpans[affectedSpans.indexOf(focusedSpan) - 1] || affectedSpans[affectedSpans.length - 1];
    onSelectSpan(preSpan);
    onSetFocusedSpan(preSpan);
  }, [affectedSpans, focusedSpan, onSelectSpan, onSetFocusedSpan]);

  return {handleNextSpan, handlePrevSpan, indexOfFocused};
};

export default useControls;
