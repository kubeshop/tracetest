import {Button, Checkbox, Typography} from 'antd';
import {ForwardOutlined} from '@ant-design/icons';
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
        Hit &apos;Skip Trace collection&apos; to create a black box test using trigger response. Handy for testing
        systems like a GET against <i>https://google.com</i> that won&apos;t send Tracetest a trace!
      </Typography.Paragraph>
      <S.Actions>
        <div>
          <Checkbox onChange={() => setShouldSave(prev => !prev)} value={shouldSave} /> Apply to all runs
        </div>
        <Button
          icon={<ForwardOutlined />}
          loading={isLoading}
          type="primary"
          ghost
          onClick={() => skipPolling(shouldSave)}
          size="small"
        >
          Skip Trace collection
        </Button>
      </S.Actions>
    </>
  );
};

export default Content;
