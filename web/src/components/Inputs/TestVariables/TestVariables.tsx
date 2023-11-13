import {useCallback} from 'react';
import {noop} from 'lodash';
import {TTestVariablesMap} from 'types/Variables.types';
import * as S from './TestVariables.styled';
import TestVariableEntry from './TestVariableEntry';

interface IProps {
  onChange?(value: Record<string, string>): void;
  value?: Record<string, string>;
  testVariables: TTestVariablesMap;
}

const TestVariables = ({onChange = noop, value = {}, testVariables}: IProps) => {
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
        <TestVariableEntry testVariables={testVariables} key={key} keyName={key} value={stringValue} onChange={handleChange} />
      ))}
    </S.TestsContainer>
  );
};

export default TestVariables;
