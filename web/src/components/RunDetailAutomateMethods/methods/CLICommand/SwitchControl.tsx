import {Switch} from 'antd';
import * as S from './CliCommand.styled';

interface IProps {
  text: string;
  value: boolean;
  id: string;
  onChange(value: boolean): void;
}

const SwitchControl = ({value, onChange, text, id}: IProps) => {
  return (
    <S.SwitchContainer>
      <Switch onChange={onChange} checked={value} id={id} />
      <S.SwitchLabel htmlFor={id}>{text}</S.SwitchLabel>
    </S.SwitchContainer>
  );
};

export default SwitchControl;
