import {noop} from 'lodash';
import {createContext, useContext, useMemo, useState} from 'react';
import ContactUsModal from './ContactUsModal';
import Launcher from './Launcher';

interface IContext {
  onOpen(): void;
}

export const Context = createContext<IContext>({
  onOpen: noop,
});

export const useContactUsModal = () => useContext(Context);

const ContactUs: React.FC = ({children}) => {
  const [isOpen, setIsOpen] = useState(false);
  const value = useMemo<IContext>(() => ({onOpen: () => setIsOpen(true)}), []);

  return (
    <Context.Provider value={value}>
      {children}
      <Launcher onClick={() => setIsOpen(true)} />
      <ContactUsModal isOpen={isOpen} onClose={() => setIsOpen(false)} />
    </Context.Provider>
  );
};

export default ContactUs;
