import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import ConfirmationModal from 'components/ConfirmationModal';

type TOnConfirm = typeof noop;
type TOnOPenProps = {
  title: React.ReactNode;
  heading?: string;
  okText?: string;
  cancelText?: string;
  onConfirm: TOnConfirm;
};

interface IContext {
  onOpen(props: TOnOPenProps): void;
}

export const Context = createContext<IContext>({
  onOpen: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useConfirmationModal = () => useContext(Context);

const ConfirmationModalProvider = ({children}: IProps) => {
  const [{title, heading, okText, cancelText, onConfirm}, setProps] = useState<TOnOPenProps>({
    title: '',
    heading: '',
    okText: '',
    cancelText: '',
    onConfirm: noop,
  });
  const [isOpen, setIsOpen] = useState(false);

  const onOpen = useCallback((newProps: TOnOPenProps) => {
    setProps(newProps);
    setIsOpen(true);
  }, []);

  const triggerConfirm = useCallback(async () => {
    await onConfirm();
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
        okText={okText}
        cancelText={cancelText}
      />
    </Context.Provider>
  );
};

export default ConfirmationModalProvider;
