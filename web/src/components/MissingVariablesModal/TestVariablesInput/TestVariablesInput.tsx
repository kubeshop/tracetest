import {TTestVariables} from 'types/Variables.types';
import * as S from '../MissingVariablesModal.styled';
import TestVariableEntry from './TestVariableEntry';

interface IProps {
  variables: TTestVariables;
  value: Record<string, string>;
  onChange(key: string, newValue: string): void;
}

const TestVariablesInput = ({
  onChange,
  variables: {
    variables: {missing},
    test,
  },
  value,
}: IProps) => {
  return (
    <S.TestContainer>
      <S.TestName>{test.name}</S.TestName>
      {missing.map(variable => (
        <TestVariableEntry key={variable.key} variable={variable} value={value} onChange={onChange} />
      ))}
    </S.TestContainer>
  );
};

export default TestVariablesInput;
