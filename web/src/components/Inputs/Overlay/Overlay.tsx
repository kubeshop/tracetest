import {useCallback, useEffect, useRef, useState} from 'react';
import {Input} from 'antd';
import useInputActions from 'hooks/useInputActions';
import {noop} from 'lodash';
import * as S from './Overlay.styled';

interface IProps {
  onChange?(value: string): void;
  value?: string;
  isDisabled?: boolean;
}

const Overlay = ({onChange = noop, value = '', isDisabled = false}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [inputValue, setInputValue] = useState(value);
  const ref = useRef(null);

  const handler = useCallback(() => {
    setIsOpen(false);
    if (inputValue) {
      onChange(inputValue);
    } else {
      setInputValue(value);
    }
  }, [inputValue, onChange, value]);

  useInputActions(ref, handler);

  useEffect(() => {
    setInputValue(value);
  }, [value]);

  useEffect(() => {
    if (ref.current && isOpen) {
      const current = ref.current as HTMLInputElement;
      current.focus();
    }
  }, [isOpen]);

  useEffect(() => {
    if (isDisabled) setIsOpen(false);
  }, [isDisabled]);

  return isOpen ? (
    <S.InputContainer ref={ref} htmlFor="overlay-input">
      <Input
        data-cy="overlay-input"
        id="overlay-input"
        onChange={event => setInputValue(event.target.value)}
        value={inputValue}
        disabled={isDisabled}
      />
    </S.InputContainer>
  ) : (
    <S.Overlay
      onClick={e => {
        if (!isDisabled) {
          e.stopPropagation();
          setIsOpen(true);
        }
      }}
      data-cy="overlay-input-overlay"
    >
      {inputValue} <S.EditIcon />
    </S.Overlay>
  );
};

export default Overlay;
