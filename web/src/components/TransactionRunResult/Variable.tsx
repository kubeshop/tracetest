import {TEnvironmentValue} from 'types/Environment.types';
import * as S from './TransactionRunResult.styled';

interface IProps {
  value: TEnvironmentValue;
}

const Variable = ({value: {key, value}}: IProps) => {
  return (
    <S.VariableContainer data-cy={`variable-card-${key}`} key={key}>
      <S.Stack>
        <S.Text opacity={0.6}>Key</S.Text>
        <S.Text fontWeight="600">{`${key} `}</S.Text>
      </S.Stack>
      <S.Stack>
        <S.Text opacity={0.6}>Value</S.Text>
        <S.Text fontWeight="600">{`${value}`}</S.Text>
      </S.Stack>
    </S.VariableContainer>
  );
};

export default Variable;
