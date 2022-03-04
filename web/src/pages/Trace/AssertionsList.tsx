import {useState} from 'react';
import {Button, List, Typography} from 'antd';
import {ISpan} from '../../types';
import CreateAssertionModal from './CreateAssertionModal';
import {useGetTestAssertionsQuery} from '../../services/TestService';

interface IProps {
  testId: string;
  targetSpan: ISpan;
}

const AssertionList = ({testId, targetSpan}: IProps) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);
  const {data: testAssertions} = useGetTestAssertionsQuery(testId);
  return (
    <div>
      <Button onClick={() => setOpenCreateAssertion(true)}>New Assertion</Button>
      <CreateAssertionModal
        key={`KEY_${targetSpan.spanId}`}
        testId={testId}
        span={targetSpan}
        open={openCreateAssertion}
        onClose={() => setOpenCreateAssertion(false)}
      />
    </div>
  );
};

export default AssertionList;
