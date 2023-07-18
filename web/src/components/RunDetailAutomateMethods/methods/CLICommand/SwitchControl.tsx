import {Switch} from 'antd';
import {noop} from 'lodash';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import * as S from './CliCommand.styled';

interface IProps {
  text: string;
  value?: boolean;
  id: string;
  disabled?: boolean;
  help?: string;
  onChange?(value: boolean): void;
}

const SwitchControl = ({value = false, onChange = noop, text, id, disabled, help}: IProps) => (
  <S.SwitchContainer>
    <Switch onChange={onChange} checked={value} id={id} disabled={disabled} />
    <S.SwitchLabel htmlFor={id} $disabled={disabled}>
      {text}
    </S.SwitchLabel>
    {!!help && <TooltipQuestion margin={6} title={help} />}
  </S.SwitchContainer>
);

export default SwitchControl;
