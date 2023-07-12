import {Axis, Orientation} from '@visx/axis';
import {Group} from '@visx/group';
import {scaleLinear} from '@visx/scale';
import without from 'lodash/without';
import {useCallback, useMemo, useState} from 'react';
import {useTheme} from 'styled-components';

import {AxisHeight, AxisOffset, NodeHeight} from 'constants/Timeline.constants';
import TimelineModel from 'models/Timeline.model';
import TimelineService from 'services/Timeline.service';
import SpanNodeFactory from './SpanNodeFactory';
import {IProps} from './Timeline';

function tickLabelProps() {
  return {
    fill: '#687492',
    fontSize: 12,
    textAnchor: 'middle',
  } as const;
}

const Visualization = ({matchedSpans, nodeType, onNodeClick, selectedSpan, spans, width = 600}: IProps) => {
  const theme = useTheme();
  const [collapsed, setCollapsed] = useState<string[]>([]);

  const nodes = useMemo(() => TimelineModel(spans, nodeType), [spans, nodeType]);
  const filteredNodes = useMemo(() => TimelineService.getFilteredNodes(nodes, collapsed), [collapsed, nodes]);
  const [min, max] = useMemo(() => TimelineService.getMinMax(nodes), [nodes]);

  const xScale = scaleLinear({
    domain: [0, max - min],
    range: [0, width - AxisOffset],
  });

  const handleOnCollapse = useCallback((id: string) => {
    setCollapsed(prevCollapsed => {
      if (prevCollapsed.includes(id)) {
        return without(prevCollapsed, id);
      }
      return [...prevCollapsed, id];
    });
  }, []);

  return width < AxisOffset ? null : (
    <svg height={nodes.length * NodeHeight + AxisHeight} width={width}>
      <Group left={AxisOffset} top={0}>
        <Axis
          label="Duration (ms)"
          labelProps={{
            fill: theme.color.textSecondary,
            fontSize: 14,
            textAnchor: 'start',
            x: -AxisOffset,
            y: -8,
          }}
          orientation={Orientation.top}
          scale={xScale}
          stroke={theme.color.textSecondary}
          tickLabelProps={tickLabelProps}
          tickStroke={theme.color.textSecondary}
          top={20}
        />
      </Group>

      <Group className={`timeline-node-${nodeType}`} left={0} top={AxisHeight}>
        {filteredNodes.map((node, index) => (
          <SpanNodeFactory
            index={index}
            indexParent={filteredNodes.findIndex(filteredNode => filteredNode.data.id === node.data.parentId)}
            isCollapsed={collapsed.includes(node.data.id)}
            isMatched={matchedSpans.includes(node.data.id)}
            isSelected={selectedSpan === node.data.id}
            key={node.data.id}
            minStartTime={min}
            node={node}
            onClick={onNodeClick}
            onCollapse={handleOnCollapse}
            xScale={xScale}
            type={node.type}
          />
        ))}
      </Group>
    </svg>
  );
};

export default Visualization;
