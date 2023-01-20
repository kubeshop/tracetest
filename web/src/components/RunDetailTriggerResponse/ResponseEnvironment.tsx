import KeyValueRow from 'components/KeyValueRow';
import {ENVIRONMENTS_DOCUMENTATION_URL} from 'constants/Common.constants';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './RunDetailTriggerResponse.styled';

const ResponseEnvironment = () => {
  const {
    run: {environment},
  } = useTestRun();

  if (!environment?.values?.length) {
    return (
      <S.EmptyContainer>
        <S.EmptyIcon />
        <S.EmptyTitle>There are no environment variables used in this test</S.EmptyTitle>
        <S.EmptyText>
          Learn more about environments{' '}
          <a href={ENVIRONMENTS_DOCUMENTATION_URL} target="_blank">
            here
          </a>
        </S.EmptyText>
      </S.EmptyContainer>
    );
  }

  return (
    <S.ResponseEnvironmentContainer>
      {environment?.values?.map(value => (
        <KeyValueRow key={value.key} keyName={value.key} value={value.value} />
      ))}
    </S.ResponseEnvironmentContainer>
  );
};

export default ResponseEnvironment;
