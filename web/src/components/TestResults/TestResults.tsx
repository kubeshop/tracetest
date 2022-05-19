import SkeletonTable from 'components/SkeletonTable';
import {FC} from 'react';
import {ITestRunResult} from '../../types/TestRunResult.types';
import AssertionCardList from '../AssertionCardList';
import {IAssertionResult} from '../../types/Assertion.types';

type TTestResultsProps = {
  assertionResultList: IAssertionResult[];
  onSelectSpan(spanId: string): void;
  result: ITestRunResult;
};

const TestResults: FC<TTestResultsProps> = ({result: {resultId, testId, trace}, onSelectSpan, assertionResultList}) => {
  return (
    <SkeletonTable loading={!trace}>
      <AssertionCardList
        assertionResultList={assertionResultList}
        onSelectSpan={onSelectSpan}
        resultId={resultId}
        testId={testId}
      />
    </SkeletonTable>
  );
};

export default TestResults;
