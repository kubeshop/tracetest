import {Button, Checkbox, Typography} from 'antd';
import {useState} from 'react';
import * as S from './SkipPollingPopover.styled';

interface IProps {
  skipPolling(shouldSave: boolean): void;
  isLoading: boolean;
}

const Content = ({skipPolling, isLoading}: IProps) => {
  const [shouldSave, setShouldSave] = useState(false);

  return (
    <>
      <Typography.Paragraph>Skip the Trace Collection step to use current state to create tests</Typography.Paragraph>
      <S.Actions>
        <div>
          <Checkbox onChange={() => setShouldSave(prev => !prev)} value={shouldSave} /> Apply to all runs
        </div>
        <Button loading={isLoading} type="primary" onClick={() => skipPolling(shouldSave)}>
          Skip
        </Button>
      </S.Actions>
    </>
  );
};

export default Content;
