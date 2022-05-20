import {noop} from 'lodash';
import {useState, createContext, useCallback, useMemo, useContext} from 'react';
import {IValues} from './AssertionForm';

interface IFormProps {
  defaultValues?: IValues;
  assertionId?: string;
  isEditing?: boolean;
}

interface ICreateAssertionModalProviderContext {
  isOpen: boolean;
  open(props?: IFormProps): void;
  close(): void;
  onSubmit(values: IValues): void;
  formProps: IFormProps;
}

const initialFormProps = {
  isEditing: false,
};

export const Context = createContext<ICreateAssertionModalProviderContext>({
  isOpen: false,
  open: noop,
  close: noop,
  formProps: initialFormProps,
  onSubmit: noop,
});

export const useAssertionForm = () => useContext(Context);

const AssertionFormProvider: React.FC<{testId: string}> = ({children}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [formProps, setFormProps] = useState<IFormProps>(initialFormProps);

  const open = useCallback((props: IFormProps = {}) => {
    setFormProps(props);
    setIsOpen(true);
  }, []);

  const close = useCallback(() => {
    setIsOpen(false);
    setFormProps(initialFormProps);
  }, []);

  const onSubmit = useCallback(({selectorList, assertionList}: IValues) => {
    console.log('@@onSubmit', selectorList, assertionList);
    setIsOpen(false);
  }, []);

  const contextValue = useMemo(
    () => ({isOpen, open, close, formProps, onSubmit}),
    [isOpen, open, close, formProps, onSubmit]
  );

  return <Context.Provider value={contextValue}>{children}</Context.Provider>;
};

export default AssertionFormProvider;
