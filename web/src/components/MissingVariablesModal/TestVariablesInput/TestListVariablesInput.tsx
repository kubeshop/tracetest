import {useCallback} from 'react';
import {noop} from 'lodash';
import {TTestVariablesMap} from 'types/Variables.types';
import TestVariablesInput from './TestVariablesInput';
import * as S from '../MissingVariablesModal.styled';

interface IProps {
  variables: TTestVariablesMap;
  onChange?(value: Record<string, string>): void;
  value?: Record<string, string>;
}

const TestListVariablesInput = ({onChange = noop, variables: variablesMap, value = {}}: IProps) => {
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
      {Object.entries(variablesMap).map(([testId, variables]) => (
        <TestVariablesInput key={testId} value={value} variables={variables} onChange={handleChange} />
      ))}
    </S.TestsContainer>
  );
};

export default TestListVariablesInput;
