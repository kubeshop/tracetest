import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import ConfirmationModal from 'components/ConfirmationModal';

type TOnConfirm = typeof noop;

interface IContext {
  onOpen(title: string, onConfirm: TOnConfirm, heading?: string): void;
}

export const Context = createContext<IContext>({
  onOpen: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useConfirmationModal = () => useContext(Context);

const ConfirmationModalProvider = ({children}: IProps) => {
  const [title, setTitle] = useState<string>('');
  const [heading, setHeading] = useState<string>('');
  const [onConfirm, setOnConfirm] = useState<TOnConfirm>(() => noop);
  const [isOpen, setIsOpen] = useState(false);

  const onOpen = useCallback((newTitle: string, onConfirmFn: TOnConfirm, newHeading = 'Delete Confirmation') => {
    setTitle(newTitle);
    setOnConfirm(() => onConfirmFn);
    setIsOpen(true);
    setHeading(newHeading);
  }, []);

  const triggerConfirm = useCallback(() => {
    onConfirm();
    setIsOpen(false);
  }, [onConfirm]);

  const value = useMemo<IContext>(() => ({onOpen}), [onOpen]);

  return (
    <Context.Provider value={value}>
      {children}
      <ConfirmationModal
        onClose={() => setIsOpen(false)}
        onConfirm={triggerConfirm}
        isOpen={isOpen}
        title={title}
        heading={heading}
      />
    </Context.Provider>
  );
};

export default ConfirmationModalProvider;
