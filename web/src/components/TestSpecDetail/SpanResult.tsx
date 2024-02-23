import {useCallback, useMemo} from 'react';
import {useSpan} from 'providers/Span/Span.provider';
import {ICheckResult} from 'types/Assertion.types';

import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import AssertionService from 'services/Assertion.service';
import * as S from './TestSpecDetail.styled';
import Detail from '../SpanResultDetail/Detail';

interface IProps {
  index: number;
  data: [string, ICheckResult[]][];
  style: React.CSSProperties;
}

const SpanResult = ({index, data, style}: IProps) => {
  const [spanId, checkResults] = useMemo(() => data[index], [data, index]);
  const {
    run: {trace},
  } = useTestRun();
  const {selectedSpan, onSetFocusedSpan, onSelectSpan} = useSpan();
  const {setSelectedSpanResult} = useTestSpecs();

  const onFocusAndSelect = useCallback(() => {
    onSelectSpan(spanId);
    onSetFocusedSpan(spanId);
    setSelectedSpanResult({spanId, checkResults});
  }, [checkResults, onSelectSpan, onSetFocusedSpan, setSelectedSpanResult, spanId]);

  const span = trace?.flat[spanId] || {};

  const assertionsPassed = useMemo(() => AssertionService.getTotalPassedSpanChecks(checkResults), [checkResults]);
  const assertionsFailed = checkResults.length - assertionsPassed;

  return (
    <S.Wrapper
      style={{
        ...style,
        bottom: Number(style.top) - 16,
        height: Number(style.height) - 16,
      }}
      $type={span.type}
      $isSelected={spanId === selectedSpan?.id}
      onClick={onFocusAndSelect}
    >
      <Detail span={span} assertionsFailed={assertionsFailed} assertionsPassed={assertionsPassed} />
    </S.Wrapper>
  );
};

export default SpanResult;
