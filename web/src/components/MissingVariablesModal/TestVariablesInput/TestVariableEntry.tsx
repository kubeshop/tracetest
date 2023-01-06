import {Input, Popover} from 'antd';
import {useMemo} from 'react';
import {Link} from 'react-router-dom';
import {TTestVariablesMap} from 'types/Variables.types';
import VariablesService from 'services/Variables.service';
import * as S from '../MissingVariablesModal.styled';

interface IProps {
  variables: TTestVariablesMap;
  keyName: string;
  value: string;
  onChange(key: string, newValue: string): void;
}

const TestVariableEntry = ({value, keyName, onChange, variables}: IProps) => {
  const toolTip = useMemo(() => {
    const isMoreThanOneTest = Object.values(variables).length > 1;
    const testList = VariablesService.getMatchingTests(variables, keyName);

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
  }, [keyName, variables]);

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
