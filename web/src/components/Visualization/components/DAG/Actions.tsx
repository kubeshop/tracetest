import {ApartmentOutlined, ExpandOutlined, ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';
import {Tooltip} from 'antd';
import {noop} from 'lodash';
import {useCallback, useEffect} from 'react';
import {NodeInternals, useReactFlow, useStoreApi} from 'react-flow-renderer';

import Navigation from '../Navigation';
import * as S from './DAG.styled';

interface IProps {
  isMiniMapActive?: boolean;
  matchedSpans: string[];
  onMiniMapToggle?: () => void;
  onNavigateToSpan(spanId: string): void;
  selectedSpan: string;
}

function getNodePosition(nodeId: string, nodeInternals: NodeInternals) {
  const nodes = Array.from(nodeInternals).map(([, node]) => node);
  const node = nodes.find(({id}) => id === nodeId);

  if (!node) return {x: 0, y: 0};

  const x = node.position.x + (node?.width ?? 0) / 2;
  const y = node.position.y + (node?.height ?? 0) / 2;

  return {x, y};
}

const Actions = ({
  isMiniMapActive = false,
  matchedSpans,
  onMiniMapToggle = noop,
  onNavigateToSpan,
  selectedSpan,
}: IProps) => {
  const {fitView, setCenter, zoomIn, zoomOut} = useReactFlow();
  const {getState} = useStoreApi();

  useEffect(() => {
    fitView({duration: 1000});
  }, [fitView, matchedSpans.length]);

  const handleOnNavigateToSpan = useCallback(
    (spanId: string) => {
      onNavigateToSpan(spanId);
    },
    [onNavigateToSpan]
  );

  useEffect(() => {
    if (selectedSpan) {
      const {nodeInternals} = getState();
      const {x, y} = getNodePosition(selectedSpan, nodeInternals);
      setCenter(x, y, {zoom: 1.0, duration: 1000});
    }
  }, [getState, selectedSpan, setCenter]);

  return (
    <>
      <Navigation matchedSpans={matchedSpans} onNavigateToSpan={handleOnNavigateToSpan} selectedSpan={selectedSpan} />

      <S.ActionsContainer>
        <Tooltip placement="right" title="Zoom In">
          <S.ActionButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
        </Tooltip>
        <Tooltip placement="right" title="Zoom Out">
          <S.ActionButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
        </Tooltip>
        <Tooltip placement="right" title="Fit View">
          <S.ActionButton icon={<ExpandOutlined />} onClick={() => fitView()} type="text" />
        </Tooltip>
        <Tooltip placement="right" title="Mini Map">
          <S.ActionButton
            icon={<ApartmentOutlined />}
            onClick={onMiniMapToggle}
            type="text"
            $isActive={isMiniMapActive}
          />
        </Tooltip>
      </S.ActionsContainer>
    </>
  );
};

export default Actions;
