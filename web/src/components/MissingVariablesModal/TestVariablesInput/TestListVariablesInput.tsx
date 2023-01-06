import {useCallback} from 'react';
import {noop} from 'lodash';
import {TTestVariablesMap} from 'types/Variables.types';
import * as S from '../MissingVariablesModal.styled';
import TestVariableEntry from './TestVariableEntry';

interface IProps {
  variables: TTestVariablesMap;
  onChange?(value: Record<string, string>): void;
  value?: Record<string, string>;
}

const TestListVariablesInput = ({onChange = noop, variables, value = {}}: IProps) => {
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
        <TestVariableEntry key={key} keyName={key} variables={variables} value={stringValue} onChange={handleChange} />
      ))}
    </S.TestsContainer>
  );
};

export default TestListVariablesInput;
