import {useNavigate} from 'react-router-dom';

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
  const navigate = useNavigate();
  const {run} = useTestRun();
  const createdTimeAgo = Date.getTimeAgo(run?.createdAt ?? '');

  return (
    <S.Section $justifyContent="flex-start">
      <S.BackIcon data-cy="test-header-back-button" onClick={() => navigate(`/test/${testId}`)} />
      <S.InfoContainer>
        <S.Row>
          <S.Title data-cy="test-details-name">
            {name} (v{run.testVersion})
          </S.Title>
          <Info
            date={run?.createdAt ?? ''}
            executionTime={run?.executionTime ?? 0}
            totalSpans={run?.trace?.spans?.length ?? 0}
            traceId={run?.traceId ?? ''}
          />
        </S.Row>
        <S.Text>
          {triggerType} • {`Ran ${createdTimeAgo}`}
        </S.Text>
      </S.InfoContainer>
    </S.Section>
  );
};

export default HeaderLeft;
