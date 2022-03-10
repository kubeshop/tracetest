import {useState} from 'react';
import {Button} from 'antd';
import {ISpan} from 'types';
import {useGetTestAssertionsQuery} from 'services/TestService';

import CreateAssertionModal from './CreateAssertionModal';

interface IProps {
  testId: string;
  targetSpan: ISpan;
  trace: any;
}

const AssertionList = ({testId, targetSpan, trace}: IProps) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);
  const {data: testAssertions} = useGetTestAssertionsQuery(testId);
  return (
    <div>
      <Button onClick={() => setOpenCreateAssertion(true)}>New Assertion</Button>
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
