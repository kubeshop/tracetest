import SkeletonTable from 'components/SkeletonTable';
import {FC} from 'react';
import { TAssertionResults } from '../../types/Assertion.types';
import AssertionCardList from '../AssertionCardList';

type TTestResultsProps = {
  onSelectSpan(spanId: string): void;
  testId: string;
  assertionResults: TAssertionResults;
};

const TestResults: FC<TTestResultsProps> = ({assertionResults, testId, onSelectSpan}) => {
  return (
    <SkeletonTable loading={!assertionResults}>
      <AssertionCardList assertionResults={assertionResults} onSelectSpan={onSelectSpan} testId={testId} />
    </SkeletonTable>
  );
};

export default TestResults;
