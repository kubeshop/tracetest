import {Button, Popover} from 'antd';
import {useMemo} from 'react';
import Link from 'components/Link';
import * as S from './NoTracingPopover.styled';

const NoTracingPopover = () => {
  const content = useMemo(
    () => (
      <S.MessageContainer>
        <div>
          <S.Title>Tracing is not configured.</S.Title>
          <S.Text>
            Tracetest needs configuration information about your Tracing Backend so it can collect the trace after a
            test run. Please configure your Tracing Backend now.
          </S.Text>
        </div>
        <S.ButtonContainer>
          <Link to="/settings">
            <Button type="primary">Configure</Button>
          </Link>
        </S.ButtonContainer>
      </S.MessageContainer>
    ),
    []
  );

  return (
    <Popover content={content} overlayClassName="no-tracing-popover" placement="bottom">
      <S.Trigger>
        <S.Icon /> No-Tracing Mode
      </S.Trigger>
    </Popover>
  );
};

export default NoTracingPopover;
