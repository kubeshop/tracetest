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
      <Typography.Paragraph>
        You can skip the <b>awaiting trace</b> step and use the current state to create test specs.
      </Typography.Paragraph>
      <S.Actions>
        <div>
          <Checkbox onChange={() => setShouldSave(prev => !prev)} value={shouldSave} /> Apply to all runs
        </div>
        <Button loading={isLoading} type="primary" ghost onClick={() => skipPolling(shouldSave)} size="small">
          Skip awaiting trace
        </Button>
      </S.Actions>
    </>
  );
};

export default Content;
