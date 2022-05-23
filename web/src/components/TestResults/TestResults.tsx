import SkeletonTable from 'components/SkeletonTable';
import {FC} from 'react';
import {TTestRun} from '../../types/TestRun.types';
import AssertionCardList from '../AssertionCardList';

type TTestResultsProps = {
  onSelectSpan(spanId: string): void;
  testId: string;
  run: TTestRun;
};

const TestResults: FC<TTestResultsProps> = ({run: {trace, result}, testId, onSelectSpan}) => {
  return (
    <SkeletonTable loading={!trace}>
      <AssertionCardList assertionResults={result} onSelectSpan={onSelectSpan} testId={testId} />
    </SkeletonTable>
  );
};

export default TestResults;
