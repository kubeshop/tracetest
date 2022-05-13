import {noop} from 'lodash';
import {useState, createContext, useCallback, useMemo, useContext} from 'react';
import {IAssertion} from '../../types/Assertion.types';
import {ISpan, ISpanFlatAttribute} from '../../types/Span.types';
import CreateAssertionModal from './CreateAssertionModal';

interface IModalProps {
  span?: ISpan;
  testId: string;
  resultId: string;
  assertion?: IAssertion;
  defaultAttributeList?: ISpanFlatAttribute[];
}

interface ICreateAssertionModalProviderContext {
  isOpen: boolean;
  open(props: IModalProps): void;
  close(): void;
}

export const Context = createContext<ICreateAssertionModalProviderContext>({
  isOpen: false,
  open: noop,
  close: noop,
});

const initialModalProps = {
  testId: '',
  resultId: '',
};

export const useCreateAssertionModal = () => useContext(Context);

const CreateAssertionModalProvider: React.FC = ({children}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [{span, ...modalProps}, setModalProps] = useState<IModalProps>(initialModalProps);

  const open = useCallback((props: IModalProps) => {
    setModalProps(props);
    setIsOpen(true);
  }, []);

  const close = useCallback(() => {
    setIsOpen(false);
    setModalProps(initialModalProps);
  }, []);

  const contextValue = useMemo(() => ({isOpen, open, close}), [isOpen, open, close]);

  return (
    <Context.Provider value={contextValue}>
      {children}
      <CreateAssertionModal open={isOpen} onClose={() => setIsOpen(false)} span={span!} {...modalProps} />
    </Context.Provider>
  );
};

export default CreateAssertionModalProvider;
