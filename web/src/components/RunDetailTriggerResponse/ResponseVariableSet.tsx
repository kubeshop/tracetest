import KeyValueRow from 'components/KeyValueRow';
import {VARIABLE_SET_DOCUMENTATION_URL} from 'constants/Common.constants';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './RunDetailTriggerResponse.styled';

const ResponseVariableSet = () => {
  const {
    run: {variableSet},
  } = useTestRun();

  if (!variableSet?.values?.length) {
    return (
      <S.EmptyContainer>
        <S.EmptyIcon />
        <S.EmptyTitle>There are no variable sets used in this test</S.EmptyTitle>
        <S.EmptyText>
          Learn more about variable sets{' '}
          <a href={VARIABLE_SET_DOCUMENTATION_URL} target="_blank">
            here
          </a>
        </S.EmptyText>
      </S.EmptyContainer>
    );
  }

  return (
    <S.ResponseVarsContainer>
      {variableSet?.values?.map(value => (
        <KeyValueRow key={value.key} keyName={value.key} value={value.value} />
      ))}
    </S.ResponseVarsContainer>
  );
};

export default ResponseVariableSet;
