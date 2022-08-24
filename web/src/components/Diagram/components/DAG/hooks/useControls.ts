import {useCallback, useEffect} from 'react';
import {useReactFlow, useStoreApi} from 'react-flow-renderer';
import {useSpan} from 'providers/Span/Span.provider';

const useControls = () => {
  const {fitView, setCenter} = useReactFlow();
  const {affectedSpans, focusedSpan, onSelectSpan, onSetFocusedSpan} = useSpan();
  const {getState} = useStoreApi();

  const indexOfFocused = affectedSpans.findIndex(spanId => spanId === focusedSpan) + 1;

  useEffect(() => {
    const getNodePosition = (nodeId: string) => {
      const {nodeInternals} = getState();
      const nodes = Array.from(nodeInternals).map(([, node]) => node);
      const node = nodes.find(({id}) => id === nodeId);

      if (!node) return {x: 0, y: 0};

      const x = node.position.x + (node?.width ?? 0) / 2;
      const y = node.position.y + (node?.height ?? 0) / 2;

      return {x, y};
    };

    if (focusedSpan) {
      const {x, y} = getNodePosition(focusedSpan);
      setCenter(x, y, {zoom: 1.0, duration: 1000});
      onSelectSpan(focusedSpan);
    } else {
      fitView({duration: 1000});
    }
  }, [fitView, focusedSpan, getState, onSelectSpan, setCenter]);

  const handleNextSpan = useCallback(() => {
    const nextSpan = affectedSpans[affectedSpans.indexOf(focusedSpan) + 1] || affectedSpans[0];
    onSetFocusedSpan(nextSpan);
  }, [affectedSpans, focusedSpan, onSetFocusedSpan]);

  const handlePrevSpan = useCallback(() => {
    const preSpan = affectedSpans[affectedSpans.indexOf(focusedSpan) - 1] || affectedSpans[affectedSpans.length - 1];
    onSetFocusedSpan(preSpan);
  }, [affectedSpans, focusedSpan, onSetFocusedSpan]);

  return {handleNextSpan, handlePrevSpan, indexOfFocused};
};

export default useControls;
