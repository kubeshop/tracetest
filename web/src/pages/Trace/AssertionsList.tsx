import {useState} from 'react';
import {Button} from 'antd';
import {ISpan} from '../../types';
import CreateAssertionModal from './CreateAssertionModal';

interface IProps {
  targetSpan: ISpan;
}

const AssertionList = ({targetSpan}: IProps) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);
  return (
    <div>
      <Button onClick={() => setOpenCreateAssertion(true)}>New Assertion</Button>
      <CreateAssertionModal
        span={targetSpan}
        open={openCreateAssertion}
        onClose={() => setOpenCreateAssertion(false)}
      />
    </div>
  );
};

export default AssertionList;
