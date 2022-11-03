import {TEnvironmentValue} from 'types/Environment.types';
import * as S from './TransactionRunResult.styled';

interface IProps {
  value: TEnvironmentValue;
}

const Variable = ({value: {key, value}}: IProps) => {
  return (
    <S.Container data-cy={`variable-card-${key}`} key={key}>
      <S.Info>
        <S.Stack>
          <S.Text opacity={0.6}>Key</S.Text>
          <S.Text>{`${key} `}</S.Text>
        </S.Stack>
        <S.Stack>
          <S.Text opacity={0.6}>Value</S.Text>
          <S.Text>{`${value}`}</S.Text>
        </S.Stack>
      </S.Info>
    </S.Container>
  );
};

export default Variable;
