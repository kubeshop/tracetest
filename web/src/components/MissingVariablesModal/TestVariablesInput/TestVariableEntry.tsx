import {Input} from 'antd';
import * as S from '../MissingVariablesModal.styled';

interface IProps {
  keyName: string;
  value: string;
  onChange(key: string, newValue: string): void;
}

const TestVariableEntry = ({value, keyName, onChange}: IProps) => {

  return (
    <S.InputContainer>
      <S.FromItem label="Value Name" $hasValue>
        <Input disabled value={keyName} />
      </S.FromItem>
      <S.FromItem label="Variable Value" $hasValue={!!value}>
        <Input
          value={value}
          placeholder="provide variable value"
          onChange={event => onChange(keyName, event.target.value)}
        />
      </S.FromItem>
    </S.InputContainer>
  );
};

export default TestVariableEntry;
