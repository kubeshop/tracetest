import {noop} from 'lodash';

import VersionMismatchModal from 'components/VersionMismatchModal/VersionMismatchModal';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setSelectedAssertion} from 'redux/slices/TestDefinition.slice';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TTestDefinitionEntry} from 'types/TestDefinition.types';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {useSpan} from 'providers/Span/Span.provider';
import {IValues} from './AssertionForm';

interface IFormProps {
  defaultValues?: IValues;
  selector?: string;
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

export const useAssertionForm = () => useContext<ICreateAssertionModalProviderContext>(Context);

const AssertionFormProvider: React.FC<{testId: string}> = ({children}) => {
  const dispatch = useAppDispatch();
  const [isOpen, setIsOpen] = useState(false);
  const [isConfirmationModalOpen, setIsConfirmationModalOpen] = useState(false);
  const [formProps, setFormProps] = useState<IFormProps>(initialFormProps);
  const {update, add, test, isDraftMode} = useTestDefinition();
  const {run} = useTestRun();
  const {onClearAffectedSpans} = useSpan();
  const definitionList = useAppSelector(state => TestDefinitionSelectors.selectDefinitionList(state));

  const open = useCallback(
    (props: IFormProps = {}) => {
      const {
        isEditing,
        defaultValues: {
          assertionList = [],
          selector: defaultSelector,
        } = {},
      } = props;
      const definition = definitionList.find(({selector}) => defaultSelector === selector);

      if (definition)
        setFormProps({
          ...props,
          isEditing: true,
          selector: defaultSelector,
          defaultValues: {
            selector: defaultSelector,
            assertionList: isEditing ? assertionList : [...definition.assertionList, ...assertionList],
          },
        });
      else setFormProps(props);

      if (run.testVersion !== test?.version && !isDraftMode) {
        CreateAssertionModalAnalyticsService.onConfirmationModalOpen();
        setIsConfirmationModalOpen(true);
      } else {
        CreateAssertionModalAnalyticsService.onAssertionFormOpen();
        setIsOpen(true);
      }

      dispatch(setSelectedAssertion());
    },
    [dispatch, definitionList, isDraftMode, run.testVersion, test?.version]
  );

  const close = useCallback(() => {
    setFormProps(initialFormProps);
    onClearAffectedSpans();

    setIsOpen(false);
  }, [onClearAffectedSpans]);

  const onConfirm = useCallback(() => {
    CreateAssertionModalAnalyticsService.onAssertionFormOpen();
    setIsOpen(true);
    setIsConfirmationModalOpen(false);
  }, []);

  const onSubmit = useCallback(
    async ({
      assertionList = [],
      selector: newSelectorString = '',
    }: IValues) => {
      const {isEditing, selector = ''} = formProps;

      const definition: TTestDefinitionEntry = {
        selector: newSelectorString,
        assertionList,
        originalSelector: newSelectorString,
        isDraft: true,
      };

      if (isEditing) {
        CreateAssertionModalAnalyticsService.onCreateAssertionFormSubmit();
        await update(selector, definition);
      } else {
        CreateAssertionModalAnalyticsService.onEditAssertionFormSubmit();
        await add(definition);
      }

      setIsOpen(false);
      onClearAffectedSpans();
    },
    [formProps, onClearAffectedSpans, update, add]
  );

  const contextValue = useMemo(
    () => ({isOpen, open, close, formProps, onSubmit}),
    [isOpen, open, close, formProps, onSubmit]
  );

  return (
    <Context.Provider value={contextValue}>
      {children}
      <VersionMismatchModal
        description="Changing and saving changes will result in a new version that will become the latest."
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
