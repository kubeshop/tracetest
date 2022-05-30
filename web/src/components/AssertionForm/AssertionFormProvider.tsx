import {noop} from 'lodash';
import {useState, createContext, useCallback, useMemo, useContext, Dispatch, SetStateAction} from 'react';

import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {clearAffectedSpans} from 'redux/slices/TestDefinition.slice';
import SelectorService from 'services/Selector.service';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TTestDefinitionEntry} from 'types/TestDefinition.types';
import {IValues} from './AssertionForm';
import AssertionFormConfirmModal from './AssertionFormConfirmModal';

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
  const dispatch = useAppDispatch();
  const [isCollapsed, setIsCollapsed] = useState(true);
  const [isOpen, setIsOpen] = useState(false);
  const [isConfirmationModalOpen, setIsConfirmationModalOpen] = useState(false);
  const [formProps, setFormProps] = useState<IFormProps>(initialFormProps);
  const {update, add, test, isDraftMode} = useTestDefinition();
  const {run} = useTestRun();
  const definitionList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionList(state));

  const open = useCallback(
    (props: IFormProps = {}) => {
      const {isEditing, defaultValues: {selectorList = [], pseudoSelector, assertionList = []} = {}} = props;
      const selectorString = SelectorService.getSelectorString(selectorList, pseudoSelector);
      const definition = definitionList.find(({selector}) => selectorString === selector);

      if (definition)
        setFormProps({
          ...props,
          isEditing: true,
          selector: selectorString,
          defaultValues: {
            pseudoSelector,
            assertionList: isEditing ? assertionList : [...definition.assertionList, ...assertionList],
            selectorList,
          },
        });
      else setFormProps(props);

      if (run.testVersion !== test?.version && !isDraftMode) setIsConfirmationModalOpen(true);
      else setIsOpen(true);
    },
    [definitionList, isDraftMode, run.testVersion, test?.version]
  );

  const close = useCallback(() => {
    setFormProps(initialFormProps);
    dispatch(clearAffectedSpans());

    setIsOpen(false);
  }, []);

  const onConfirm = useCallback(() => {
    setIsOpen(true);
    setIsConfirmationModalOpen(false);
  }, []);

  const onSubmit = useCallback(
    async ({selectorList, assertionList = [], pseudoSelector}: IValues) => {
      const {isEditing, selector = ''} = formProps;

      const definition: TTestDefinitionEntry = {
        selector: SelectorService.getSelectorString(selectorList, pseudoSelector),
        assertionList,
        isDraft: true,
      };

      if (isEditing) await update(selector, definition);
      else await add(definition);

      setIsOpen(false);
      dispatch(clearAffectedSpans());
    },
    [add, formProps, update]
  );

  const contextValue = useMemo(
    () => ({isOpen, open, close, formProps, onSubmit, isCollapsed, setIsCollapsed}),
    [isOpen, open, close, formProps, onSubmit, isCollapsed, setIsCollapsed]
  );

  return (
    <Context.Provider value={contextValue}>
      {children}
      <AssertionFormConfirmModal
        isOpen={isConfirmationModalOpen}
        latestVersion={test?.version || 1}
        currentVersion={run.testVersion}
        onCancel={() => setIsConfirmationModalOpen(false)}
        onConfirm={onConfirm}
      />
    </Context.Provider>
  );
};

export default AssertionFormProvider;
