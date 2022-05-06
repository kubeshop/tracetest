import React, {Dispatch, PointerEventHandler, SetStateAction, useEffect} from 'react';
import Title from 'antd/lib/typography/Title';

import SkeletonTable from 'components/SkeletonTable';
import {Steps} from 'components/GuidedTour/traceStepList';
import {ITrace} from 'types/Trace.types';
import {ISpan} from 'types/Span.types';
import {TimelineChart} from './TimelineChart';
import {TraceHeader} from './TestResults.styled';
import {useElementSize} from '../../../hooks/useElementSize';
import {useDoubleClick} from '../../../hooks/useDoubleClick';
import GuidedTourService, {GuidedTours} from '../../../services/GuidedTour.service';
import './TimelineChart.css';

interface IProps {
  onPointerDown?: PointerEventHandler;
  trace?: ITrace;
  visiblePortion: number;
  setMax: Dispatch<SetStateAction<number>>;
  max: number;
  height?: number;
  setHeight?: Dispatch<SetStateAction<number>>;
  selectedSpan?: ISpan;

  onSelectSpan(spanId: string): void;
}

export const TraceTimeline = ({
  setHeight,
  setMax,
  onPointerDown,
  visiblePortion,
  trace,
  selectedSpan,
  onSelectSpan,
  ...props
}: IProps) => {
  const [squareRef, {height}] = useElementSize();
  useEffect(() => setMax(height), [height, setMax]);
  return (
    <div ref={squareRef} onClick={useDoubleClick(() => setHeight?.(height === props.height ? visiblePortion : height))}>
      <TraceHeader
        onPointerDown={onPointerDown}
        visiblePortion={visiblePortion}
        data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Timeline)}
        style={{height: visiblePortion}}
      >
        <Title level={4} style={{margin: 0}}>
          Component Timeline
        </Title>
      </TraceHeader>
      <SkeletonTable loading={!trace || !selectedSpan}>
        <TimelineChart trace={trace!} selectedSpan={selectedSpan} onSelectSpan={onSelectSpan} />
      </SkeletonTable>
    </div>
  );
};
