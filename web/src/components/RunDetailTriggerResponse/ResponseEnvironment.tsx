import KeyValueRow from 'components/KeyValueRow';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './RunDetailTriggerResponse.styled';

const ResponseEnvironment = () => {
  const {
    run: {environment},
  } = useTestRun();

  return (
    <S.ResponseEnvironmentContainer>
      {environment?.values?.map(value => (
        <KeyValueRow key={value.key} keyName={value.key} value={value.value} />
      ))}
    </S.ResponseEnvironmentContainer>
  );
};

export default ResponseEnvironment;
