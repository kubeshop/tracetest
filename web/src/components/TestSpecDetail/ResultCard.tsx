import {useCallback, useMemo} from 'react';
import {useSpan} from 'providers/Span/Span.provider';
import {useTest} from 'providers/Test/Test.provider';
import {ICheckResult} from 'types/Assertion.types';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './TestSpecDetail.styled';
import Assertion from './Assertion';
import SpanHeader from './SpanHeader';

interface IProps {
  index: number;
  data: [string, ICheckResult[]][];
  style: React.CSSProperties;
}

const ResultCard = ({index, data, style}: IProps) => {
  const [spanId, checkResults] = useMemo(() => data[index], [data, index]);
  const {
    run: {trace, id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const {selectedSpan, onSetFocusedSpan, onSelectSpan} = useSpan();

  const onFocusAndSelect = useCallback(() => {
    onSelectSpan(spanId);
    onSetFocusedSpan(spanId);
  }, [onSelectSpan, onSetFocusedSpan, spanId]);

  const span = trace?.flat[spanId];

  return (
    <S.CardContainer
      key={spanId}
      style={{
        ...style,
        bottom: Number(style.top) - 16,
        height: Number(style.height) - 16,
      }}
      size="small"
      title={<SpanHeader onSelectSpan={onFocusAndSelect} span={span} />}
      type="inner"
      $isSelected={spanId === selectedSpan?.id}
      $type={span?.type ?? SemanticGroupNames.General}
      id={`assertion-result-${spanId}`}
      onClick={() => onSelectSpan(span?.id ?? '')}
    >
      <S.AssertionsContainer>
        {checkResults.map(checkResult => (
          <Assertion
            testId={testId}
            runId={runId}
            selector=""
            check={checkResult}
            key={`${checkResult.result.spanId}-${checkResult.assertion}`}
          />
        ))}
      </S.AssertionsContainer>
    </S.CardContainer>
  );
};

export default ResultCard;
