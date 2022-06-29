import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {TTest} from 'types/Test.types';
import EditTestModal from './EditTestModal';

interface IContext {
  open(test: TTest): void;
}

export const Context = createContext<IContext>({
  open: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useEditTestModal = () => useContext(Context);

const EditTestModalProvider = ({children}: IProps) => {
  const [isOpen, setIsOpen] = useState(false);
  const [test, setTest] = useState<TTest>();

  const open = useCallback((editingTest: TTest) => {
    setTest(editingTest);
    setIsOpen(true);
  }, []);

  const value: IContext = useMemo(() => ({open}), [open]);

  return (
    <>
      <Context.Provider value={value}>{children}</Context.Provider>
      {isOpen && test && <EditTestModal isOpen test={test} onClose={() => setIsOpen(false)} />}
    </>
  );
};

export default EditTestModalProvider;
