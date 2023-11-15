import {LinkOutlined} from '@ant-design/icons';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useMemo} from 'react';
import Date from 'utils/Date';
import HeaderForm from './HeaderForm';
import Info from './Info';
import * as S from './RunDetailLayout.styled';

interface IProps {
  name: string;
  triggerType: string;
  origin: string;
}

const HeaderLeft = ({name, triggerType, origin}: IProps) => {
  const {
    run: {id: runId, createdAt, testSuiteId, testSuiteRunId, executionTime, trace, traceId, testVersion} = {},
    run,
  } = useTestRun();
  const {onEditAndReRun} = useTest();
  const createdTimeAgo = Date.getTimeAgo(createdAt ?? '');
  const {navigate} = useDashboard();

  const description = useMemo(() => {
    return (
      <>
        v{testVersion} • {triggerType} • Ran {createdTimeAgo}
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
  }, [triggerType, createdTimeAgo, testSuiteId, testSuiteRunId, testVersion]);

  return (
    <S.Section $justifyContent="flex-start">
      <a data-cy="test-header-back-button" onClick={() => navigate(origin)}>
        <S.BackIcon />
      </a>
      <S.InfoContainer>
        <S.Row>
          <HeaderForm
            name={name}
            onSubmit={draft => {
              onEditAndReRun(draft, runId ?? 1);
            }}
          />
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
