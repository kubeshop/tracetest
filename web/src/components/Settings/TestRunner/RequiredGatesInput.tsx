import {Switch, Typography} from 'antd';
import {useCallback} from 'react';
import {noop} from 'lodash';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import {SupportedRequiredGates, SupportedRequiredGatesDescription} from 'models/TestRunner.model';
import {ToTitle} from 'utils/Common';
import * as S from '../common/Settings.styled';

interface IProps {
  value?: string[];
  onChange?(value: string[]): void;
}

const supportedGates = Object.values(SupportedRequiredGates);

const RequiredGatesInput = ({value = [], onChange = noop}: IProps) => {
  const handleChange = useCallback(
    (gate: SupportedRequiredGates, isChecked: boolean) => {
      const newValue = isChecked ? [...value, gate] : value.filter(g => g !== gate);
      onChange(newValue);
    },
    [onChange, value]
  );

  return (
    <>
      <Typography.Title level={3}>Required Gates</Typography.Title>
      <S.SwitchListContainer>
        {supportedGates.map(gate => (
          <S.SwitchContainer key={gate}>
            <Switch checked={value.includes(gate)} onChange={isChecked => handleChange(gate, isChecked)} id={gate} />
            <label htmlFor={gate}>
              <Typography.Text>
                {ToTitle(gate)} <TooltipQuestion margin={6} title={SupportedRequiredGatesDescription[gate]} />
              </Typography.Text>
            </label>
          </S.SwitchContainer>
        ))}
      </S.SwitchListContainer>
    </>
  );
};

export default RequiredGatesInput;
