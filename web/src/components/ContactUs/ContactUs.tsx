import {PropsWithChildren, createContext, useContext, useMemo, useState} from 'react';
import {noop} from 'lodash';
import * as S from './ContactUs.styled';
import ContactUsModal from './ContactUsModal';
import PulseButton from '../PulseButton';

interface IContext {
  onOpen(): void;
}

export const Context = createContext<IContext>({
  onOpen: noop,
});

export const useContactUsModal = () => useContext(Context);

const ContactUs: React.FC<PropsWithChildren<{}>> = ({children}) => {
  const [isOpen, setIsOpen] = useState(false);
  const value = useMemo<IContext>(() => ({onOpen: () => setIsOpen(true)}), []);

  return (
    <Context.Provider value={value}>
      {children}
      <S.Container onClick={() => setIsOpen(true)}>
        <S.PulseButtonContainer>
          <PulseButton />
        </S.PulseButtonContainer>
        <S.PlushieImage />
      </S.Container>
      <ContactUsModal isOpen={isOpen} onClose={() => setIsOpen(false)} />
    </Context.Provider>
  );
};

export default ContactUs;
