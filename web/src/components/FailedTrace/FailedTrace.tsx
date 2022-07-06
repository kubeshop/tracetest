import {Button, Typography} from 'antd';
import {useCallback} from 'react';
import {DISCORD_URL, GITHUB_ISSUES_URL} from '../../constants/Common.constants';
import {useRunTestMutation} from '../../redux/apis/TraceTest.api';
import {TTestRun} from '../../types/TestRun.types';
import * as S from './FailedTrace.styled';

interface IFailedTraceProps {
  isDisplayingError: boolean;
  run: TTestRun;
  testId: string;
  onEdit(): void;
  onRunTest(result: TTestRun): void;
}

const FailedTrace: React.FC<IFailedTraceProps> = ({
  onRunTest,
  onEdit,
  testId,
  isDisplayingError,
  run: {lastErrorState},
}) => {
  const [runNewTest] = useRunTestMutation();

  const onReRun = useCallback(async () => {
    const result = await runNewTest({testId}).unwrap();
    onRunTest(result);
  }, [onRunTest, runNewTest, testId]);

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
          <Button type="primary" ghost onClick={onEdit}>
            Edit Test
          </Button>
          <Button type="primary" ghost onClick={onReRun}>
            Rerun Test
          </Button>
        </S.ButtonContainer>
      </S.Container>
    </S.FailedTrace>
  ) : null;
};

export default FailedTrace;
