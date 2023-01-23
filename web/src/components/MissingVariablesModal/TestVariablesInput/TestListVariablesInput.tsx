import {useCallback} from 'react';
import {noop} from 'lodash';
import * as S from '../MissingVariablesModal.styled';
import TestVariableEntry from './TestVariableEntry';

interface IProps {
  onChange?(value: Record<string, string>): void;
  value?: Record<string, string>;
}

const TestListVariablesInput = ({onChange = noop, value = {}}: IProps) => {
  const handleChange = useCallback(
    (key: string, newValue: string) => {
      onChange({
        ...value,
        [key]: newValue,
      });
    },
    [onChange, value]
  );

  return (
    <S.TestsContainer>
      {Object.entries(value).map(([key, stringValue]) => (
        <TestVariableEntry key={key} keyName={key} value={stringValue} onChange={handleChange} />
      ))}
    </S.TestsContainer>
  );
};

export default TestListVariablesInput;
