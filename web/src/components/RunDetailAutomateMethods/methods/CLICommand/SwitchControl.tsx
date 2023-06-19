import {Switch} from 'antd';
import * as S from './CliCommand.styled';

interface IProps {
  text: string;
  value: boolean;
  onChange(value: boolean): void;
}

const SwitchControl = ({value, onChange, text}: IProps) => {
  return (
    <S.SwitchContainer>
      <Switch onChange={onChange} checked={value} />
      {text}
    </S.SwitchContainer>
  );
};

export default SwitchControl;
