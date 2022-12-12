import {Popover} from 'antd';
import {useMemo} from 'react';
import {Link} from 'react-router-dom';
import * as S from './NoTracingPopover.styled';

const NoTracingPopover = () => {
  const content = useMemo(
    () => (
      <S.MessageContainer>
        <S.Title>
          <S.Icon /> No-Tracing Mode
        </S.Title>
        <S.Text>Tracetest is not configured to work with traces. Go to the settings page to set it up.</S.Text>
        <S.ButtonContainer>
          <Link to="/settings">
            <S.WarningButton>Setup</S.WarningButton>
          </Link>
        </S.ButtonContainer>
      </S.MessageContainer>
    ),
    []
  );

  return (
    <>
      <S.CustomPopoverGlobalStyles />
      <Popover content={content} overlayClassName="no-tracing-popover" placement="bottomLeft">
        <S.Trigger>
          <S.Icon /> No-Tracing Mode
        </S.Trigger>
      </Popover>
    </>
  );
};

export default NoTracingPopover;
