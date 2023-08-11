import {LinkOutlined} from '@ant-design/icons';
import {useMemo} from 'react';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import Date from 'utils/Date';
import Info from './Info';
import * as S from './RunDetailLayout.styled';

interface IProps {
  name: string;
  triggerType: string;
}

const HeaderLeft = ({name, triggerType}: IProps) => {
  const {run: {createdAt, testSuiteId, testSuiteRunId, executionTime, trace, traceId, testVersion} = {}, run} =
    useTestRun();
  const createdTimeAgo = Date.getTimeAgo(createdAt ?? '');
  const {navigate} = useDashboard();

  const description = useMemo(() => {
    return (
      <>
        {triggerType} • Ran {createdTimeAgo}
        {testSuiteId && testSuiteRunId && (
          <>
            {' '}
            •{' '}
            <S.TransactionLink to={`/testsuite/${testSuiteId}/run/${testSuiteRunId}`} target="_blank">
              Part of test suite <LinkOutlined />
            </S.TransactionLink>
          </>
        )}
      </>
    );
  }, [triggerType, createdTimeAgo, testSuiteId, testSuiteRunId]);

  return (
    <S.Section $justifyContent="flex-start">
      <a data-cy="test-header-back-button" onClick={() => navigate(-1)}>
        <S.BackIcon />
      </a>
      <S.InfoContainer>
        <S.Row>
          <S.Title data-cy="test-details-name">
            {name} (v{testVersion})
          </S.Title>
          <Info
            date={createdAt ?? ''}
            executionTime={executionTime ?? 0}
            state={run.state}
            totalSpans={trace?.spans?.length ?? 0}
            traceId={traceId ?? ''}
          />
        </S.Row>
        <S.Text>{description}</S.Text>
      </S.InfoContainer>
    </S.Section>
  );
};

export default HeaderLeft;
