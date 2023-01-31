import {Link} from 'react-router-dom';
import {useMemo} from 'react';
import {LinkOutlined} from '@ant-design/icons';

import {useTestRun} from 'providers/TestRun/TestRun.provider';
import Date from 'utils/Date';
import Info from './Info';
import * as S from './RunDetailLayout.styled';

interface IProps {
  name: string;
  testId: string;
  triggerType: string;
}

const HeaderLeft = ({name, testId, triggerType}: IProps) => {
  const {run: {createdAt, transactionId, transactionRunId, executionTime, trace, traceId, testVersion} = {}, run} =
    useTestRun();
  const createdTimeAgo = Date.getTimeAgo(createdAt ?? '');

  const description = useMemo(() => {
    return (
      <>
        {triggerType} • Ran {createdTimeAgo}
        {transactionId && transactionRunId && (
          <>
            {' '}
            •{' '}
            <S.TransactionLink to={`/transaction/${transactionId}/run/${transactionRunId}`} target="_blank">
              Part of transaction <LinkOutlined />
            </S.TransactionLink>
          </>
        )}
      </>
    );
  }, [createdTimeAgo, transactionId, transactionRunId, triggerType]);

  return (
    <S.Section $justifyContent="flex-start">
      <Link data-cy="test-header-back-button" to={`/test/${testId}`}>
        <S.BackIcon />
      </Link>
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
