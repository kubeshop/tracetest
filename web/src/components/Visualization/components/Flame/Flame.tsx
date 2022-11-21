import {convertJaegerTraceToProfile, FlamegraphRenderer} from '@pyroscope/flamegraph';
import '@pyroscope/flamegraph/dist/index.css';
import {useMemo} from 'react';
import FlameGraphService from 'services/Flamegraph.service';
import {TSpan} from 'types/Span.types';
import {useTestRun} from '../../../../providers/TestRun/TestRun.provider';
import * as S from './Flame.styled';

interface IProps {
  isMatchedMode: boolean;
  isTrace?: boolean;
  matchedSpans: string[];
  selectedSpan: string;
  spans: TSpan[];
  width?: number;
  onNavigateToSpan(spanId: string): void;
  onNodeClick(spanId: string): void;
}

export const Flame: React.FC<IProps> = ({spans, isTrace = false}) => {
  const {run} = useTestRun();
  const profile = useMemo(
    () => convertJaegerTraceToProfile(FlameGraphService.convertTracetestSpanToJaeger(run.trace?.traceId || '', spans)),
    [spans, run.trace?.traceId]
  );
  return (
    <S.Container $isTrace={isTrace}>
      <FlamegraphRenderer profile={profile} colorMode="light" showToolbar />
    </S.Container>
  );
};
