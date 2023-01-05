import {Input} from 'antd';
import {TMissingVariables} from 'types/Variables.types';
import * as S from '../MissingVariablesModal.styled';

interface IProps {
  variable: TMissingVariables;
  value: Record<string, string>;
  onChange(key: string, newValue: string): void;
}

const TestVariableEntry = ({value, variable: {key}, onChange}: IProps) => {
  const strValue = value[key];
  return (
    <S.InputContainer>
      <S.FromItem label="Variable Name" $hasValue>
        <Input disabled value={key} />
      </S.FromItem>
      <S.FromItem label="Variable Value" $hasValue={!!strValue}>
        <Input
          value={strValue}
          placeholder="provide variable value"
          onChange={event => onChange(key, event.target.value)}
        />
      </S.FromItem>
    </S.InputContainer>
  );
};

export default TestVariableEntry;
