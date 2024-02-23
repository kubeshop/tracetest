import {useMemo} from 'react';
import {useTest} from 'providers/Test/Test.provider';
import {ICheckResult} from 'types/Assertion.types';
import AssertionService from 'services/Assertion.service';
import Span from 'models/Span.model';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './SpanResultDetail.styled';
import Assertion from './Assertion';
import Header from './Header';
import {useTestSpecs} from '../../providers/TestSpecs/TestSpecs.provider';

interface IProps {
  span: Span;
  checkResults: ICheckResult[];
  onClose(): void;
}

const Content = ({checkResults, span, onClose}: IProps) => {
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();

  const totalPassedChecks = useMemo(() => AssertionService.getTotalPassedSpanChecks(checkResults), [checkResults]);
  const {selectedTestSpec} = useTestSpecs();

  return (
    <>
      <Header
        span={span}
        onClose={onClose}
        assertionsFailed={checkResults.length - totalPassedChecks}
        assertionsPassed={totalPassedChecks}
      />
      <S.AssertionsContainer>
        {checkResults.map(checkResult => (
          <Assertion
            testId={testId}
            runId={runId}
            selector={selectedTestSpec?.selector || ''}
            check={checkResult}
            key={`${checkResult.result.spanId}-${checkResult.assertion}`}
          />
        ))}
      </S.AssertionsContainer>
    </>
  );
};

export default Content;
