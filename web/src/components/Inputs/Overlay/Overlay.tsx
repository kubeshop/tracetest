import {useState} from 'react';
import {Input} from 'antd';
import {noop} from 'lodash';
import * as S from './Overlay.styled';

interface IProps {
  onChange?(value: string): void;
  value?: string;
}

const Overlay = ({onChange = noop, value = ''}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState(value);

  return isOpen ? (
    <S.InputContainer>
      <Input onChange={event => setInputValue(event.target.value)} value={inputValue} />
      <S.SaveButton
        ghost
        type="primary"
        onClick={() => {
          setIsOpen(false);
          onChange(inputValue);
        }}
      >
        Save
      </S.SaveButton>
    </S.InputContainer>
  ) : (
    <S.Overlay
      onClick={() => setIsOpen(true)}
    >
      {inputValue} <S.EditIcon />
    </S.Overlay>
  );
};

export default Overlay;
