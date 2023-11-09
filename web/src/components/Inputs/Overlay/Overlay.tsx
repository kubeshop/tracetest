import {useEffect, useRef, useState} from 'react';
import {Input} from 'antd';
import useInputActions from 'hooks/useInputActions';
import {noop} from 'lodash';
import * as S from './Overlay.styled';

interface IProps {
  onChange?(value: string): void;
  value?: string;
}

const Overlay = ({onChange = noop, value = ''}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState(value);
  const ref = useRef(null);
  useInputActions(ref, () => {
    onChange(inputValue);
    setIsOpen(false);
  });

  useEffect(() => {
    setInputValue(value);
  }, [value]);

  useEffect(() => {
    if (ref.current && isOpen) {
      const current = ref.current as HTMLInputElement;
      current.focus();
    }
  }, [isOpen]);

  return isOpen ? (
    <S.InputContainer ref={ref} htmlFor="overlay-input">
      <Input
        data-cy="overlay-input"
        id="overlay-input"
        onChange={event => setInputValue(event.target.value)}
        value={inputValue}
      />
    </S.InputContainer>
  ) : (
    <S.Overlay
      onClick={e => {
        e.stopPropagation();
        setIsOpen(true);
      }}
      data-cy="overlay-input-overlay"
    >
      {inputValue} <S.EditIcon />
    </S.Overlay>
  );
};

export default Overlay;
