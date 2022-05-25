import {noop} from 'lodash';
import {useState, createContext, useCallback, useMemo, useContext, Dispatch, SetStateAction} from 'react';
import {useTestDefinition} from '../../providers/TestDefinition/TestDefinition.provider';
import SelectorService from '../../services/Selector.service';
import {TTestDefinitionEntry} from '../../types/TestDefinition.types';
import {IValues} from './AssertionForm';

interface IFormProps {
  defaultValues?: IValues;
  selector?: string;
  isEditing?: boolean;
}

interface ICreateAssertionModalProviderContext {
  isCollapsed: boolean;

  setIsCollapsed: Dispatch<SetStateAction<boolean>>;

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
  isCollapsed: false,
  setIsCollapsed: noop,
  isOpen: false,
  open: noop,
  close: noop,
  formProps: initialFormProps,
  onSubmit: noop,
});

export const useAssertionForm = () => useContext(Context);

const AssertionFormProvider: React.FC<{testId: string}> = ({children}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [formProps, setFormProps] = useState<IFormProps>(initialFormProps);
  const {update, add} = useTestDefinition();

  const open = useCallback((props: IFormProps = {}) => {
    setFormProps(props);
    setIsOpen(true);
  }, []);

  const close = useCallback(() => {
    setIsOpen(false);
    setFormProps(initialFormProps);
  }, []);

  const onSubmit = useCallback(
    async ({selectorList, assertionList, pseudoSelector}: IValues) => {
      const {isEditing, selector = ''} = formProps;

      const definition: TTestDefinitionEntry = {
        selector: SelectorService.getSelectorString(selectorList, pseudoSelector),
        assertionList,
      };

      if (isEditing) await update(selector, definition);
      else await add(definition);

      setIsOpen(false);
    },
    [add, formProps, update]
  );

  const contextValue = useMemo(
    () => ({isOpen, open, close, formProps, onSubmit, isCollapsed, setIsCollapsed}),
    [isOpen, open, close, formProps, onSubmit, isCollapsed, setIsCollapsed]
  );

  return <Context.Provider value={contextValue}>{children}</Context.Provider>;
};

export default AssertionFormProvider;
