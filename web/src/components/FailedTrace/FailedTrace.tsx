import {Button, Typography} from 'antd';
import {useCallback} from 'react';
import {Link, useNavigate} from 'react-router-dom';
import {DISCORD_URL, GITHUB_ISSUES_URL} from '../../constants/Common.constants';
import {useReRunMutation} from '../../redux/apis/TraceTest.api';
import {TTestRun} from '../../types/TestRun.types';
import * as S from './FailedTrace.styled';

interface IFailedTraceProps {
  isDisplayingError: boolean;
  run: TTestRun;
  testId: string;
}

const FailedTrace: React.FC<IFailedTraceProps> = ({testId, isDisplayingError, run: {lastErrorState, id}}) => {
  const [reRunTest] = useReRunMutation();
  const navigate = useNavigate();

  const onReRun = useCallback(async () => {
    const result = await reRunTest({testId, runId: id}).unwrap();

    navigate(`/test/${testId}/run/${result.id}`);
  }, [id, navigate, reRunTest, testId]);

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
        <S.ButtonContainer>
          <Link to={`/test/${testId}/edit`}>
            <Button type="primary" ghost>
              Edit Test
            </Button>
          </Link>
          <Button type="primary" ghost onClick={onReRun}>
            Rerun Test
          </Button>
        </S.ButtonContainer>
      </S.Container>
    </S.FailedTrace>
  ) : null;
};

export default FailedTrace;
