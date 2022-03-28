import React, {useState} from 'react';
import {Button} from 'antd';

import {AssertionResult, ISpan} from 'types';
import {runTestAssertion} from 'services/AssertionService';

import {useGetTestByIdQuery} from 'services/TestService';
import AssertionsResultTable from 'components/AssertionsTable/AssertionsTable';

import CreateAssertionModal from './CreateAssertionModal';

interface IProps {
  testId: string;
  targetSpan: ISpan;
  trace: any;
}

const AssertionList = ({testId, targetSpan, trace}: IProps) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);
  const {data: test} = useGetTestByIdQuery(testId);

  const assertionsResults = test?.assertions
    ?.map(el => runTestAssertion(trace, el))
    .flat()
    .filter((f?: AssertionResult): f is AssertionResult => Boolean(f));

  return (
    <div>
      <Button style={{marginBottom: 8}} onClick={() => setOpenCreateAssertion(true)}>
        New Assertion
      </Button>
      {assertionsResults && assertionsResults.length > 0 && (
        <AssertionsResultTable assertionResults={assertionsResults} />
      )}
      <CreateAssertionModal
        key={`KEY_${targetSpan.spanId}`}
        testId={testId}
        trace={trace}
        span={targetSpan}
        open={openCreateAssertion}
        onClose={() => setOpenCreateAssertion(false)}
      />
    </div>
  );
};

export default AssertionList;
