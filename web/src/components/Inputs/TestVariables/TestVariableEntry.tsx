import {Input, Popover} from 'antd';
import {useMemo} from 'react';
import Link from 'components/Link';
import {TTestVariablesMap} from 'types/Variables.types';
import VariablesService from 'services/Variables.service';
import * as S from './TestVariables.styled';

interface IProps {
  keyName: string;
  value: string;
  onChange(key: string, newValue: string): void;
  testVariables: TTestVariablesMap;
}

const TestVariableEntry = ({value, keyName, onChange, testVariables}: IProps) => {
  const toolTip = useMemo(() => {
    const isMoreThanOneTest = Object.values(testVariables).length > 1;
    const testList = VariablesService.getMatchingTests(testVariables, keyName);

    const content = (
      <S.DetailContainer>
        <S.TestList>
          {testList.map(test => (
            <li key={test.id}>
              <Link to={`/test/${test.id}`} target="_blank">
                {test.name}
              </Link>
            </li>
          ))}
        </S.TestList>
      </S.DetailContainer>
    );

    if (isMoreThanOneTest)
      return (
        <Popover content={content} title={<S.ToolTipTitle>Used in tests</S.ToolTipTitle>}>
          <S.InfoIcon />
        </Popover>
      );

    return null;
  }, [keyName, testVariables]);

  return (
    <S.InputContainer>
      <S.FromItem label={<div>Value Name {toolTip}</div>} $hasValue>
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
