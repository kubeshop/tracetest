import {useCallback, useEffect} from 'react';
import {useReactFlow, useStoreApi} from 'react-flow-renderer';
import {useSpan} from 'providers/Span/Span.provider';

interface IProps {
  onSelectSpan(spanId: string): void;
}

const useControls = ({onSelectSpan}: IProps) => {
  const {setCenter, fitView} = useReactFlow();
  const {focusedSpan, onSetFocusedSpan, affectedSpans} = useSpan();
  const {getState} = useStoreApi();
  const indexOfFocused = affectedSpans.findIndex(spanId => spanId === focusedSpan) + 1;

  const getNodePosition = useCallback(
    (nodeId: string) => {
      const {nodeInternals} = getState();
      const nodes = Array.from(nodeInternals).map(([, node]) => node);
      const node = nodes.find(({id}) => id === nodeId);

      if (!node) return {x: 0, y: 0};

      const x = node.position.x + (node?.width ?? 0) / 2;
      const y = node.position.y + (node?.height ?? 0) / 2;

      return {x, y};
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

      setCenter(x, y, {zoom: 1.85, duration: 1000});
    } else fitView({duration: 1000});
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
