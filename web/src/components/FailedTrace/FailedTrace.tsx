import {Typography} from 'antd';
import {DISCORD_URL, GITHUB_ISSUES_URL} from 'constants/Common.constants';
import {TTestRun} from 'types/TestRun.types';
import * as S from './FailedTrace.styled';

interface IProps {
  isDisplayingError: boolean;
  run: TTestRun;
}

const FailedTrace = ({isDisplayingError, run: {lastErrorState}}: IProps) => {
  return isDisplayingError ? (
    <S.FailedTrace>
      <S.Container>
        <S.FailedIcon />
        <S.TextContainer>
          <Typography.Title level={1}>Test Run Failed</Typography.Title>
          <Typography.Text type="secondary">{lastErrorState}</Typography.Text>
          <Typography.Text type="secondary">
            Please let us know about this issue - <a href={GITHUB_ISSUES_URL}>create an issue</a> or contact us via{' '}
            <a href={DISCORD_URL}>Discord</a>.
          </Typography.Text>
          <Typography.Text type="secondary">We will check it out and respond to you.</Typography.Text>
        </S.TextContainer>
      </S.Container>
    </S.FailedTrace>
  ) : null;
};

export default FailedTrace;
