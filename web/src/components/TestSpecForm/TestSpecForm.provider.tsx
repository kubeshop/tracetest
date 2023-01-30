import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';

import VersionMismatchModal from 'components/VersionMismatchModal/VersionMismatchModal';
import {RouterSearchFields} from 'constants/Common.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useTest} from 'providers/Test/Test.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import RouterActions from 'redux/actions/Router.actions';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import AssertionService from 'services/Assertion.service';
import {TTestSpecEntry} from 'types/TestSpecs.types';
import {IValues} from './TestSpecForm';

interface IFormProps {
  defaultValues?: IValues;
  selector?: string;
  isEditing?: boolean;
}

interface IContext {
  isOpen: boolean;
  open(props?: IFormProps): void;
  close(): void;
  onSubmit(values: IValues, spanId?: string): void;
  formProps: IFormProps;
}

const initialFormProps = {
  isEditing: false,
};

export const Context = createContext<IContext>({
  isOpen: false,
  open: noop,
  close: noop,
  formProps: initialFormProps,
  onSubmit: noop,
});

export const useTestSpecForm = () => useContext<IContext>(Context);

const TestSpecFormProvider: React.FC<{testId: string}> = ({children}) => {
  const dispatch = useAppDispatch();
  const [isOpen, setIsOpen] = useState(false);
  const [isConfirmationModalOpen, setIsConfirmationModalOpen] = useState(false);
  const [formProps, setFormProps] = useState<IFormProps>(initialFormProps);
  const {test} = useTest();
  const {update, add, isDraftMode} = useTestSpecs();
  const {run} = useTestRun();
  const {onClearMatchedSpans} = useSpan();
  const specs = useAppSelector(state => TestSpecsSelectors.selectSpecs(state));

  const open = useCallback(
    (props: IFormProps = {}) => {
      const {isEditing, defaultValues: {assertions = [], selector: defaultSelector, name = ''} = {}} = props;
      const spec = specs.find(({selector}) => defaultSelector === selector);

      if (spec)
        setFormProps({
          ...props,
          isEditing: true,
          selector: defaultSelector,
          defaultValues: {
            selector: defaultSelector,
            name,
            assertions: isEditing
              ? assertions
              : [
                  ...spec.assertions.map(assertion => AssertionService.getStructuredAssertion(assertion)),
                  ...assertions,
                ],
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

      dispatch(RouterActions.updateSearch({[RouterSearchFields.SelectedAssertion]: undefined}));
    },
    [dispatch, specs, isDraftMode, run.testVersion, test?.version]
  );

  const close = useCallback(() => {
    setFormProps(initialFormProps);
    onClearMatchedSpans();

    setIsOpen(false);
  }, [onClearMatchedSpans]);

  const onConfirm = useCallback(() => {
    CreateAssertionModalAnalyticsService.onAssertionFormOpen();
    setIsOpen(true);
    setIsConfirmationModalOpen(false);
  }, []);

  const onSubmit = useCallback(
    async ({assertions = [], selector: newSelectorString = '', name = ''}: IValues) => {
      const {isEditing, selector = ''} = formProps;

      const definition: TTestSpecEntry = {
        selector: newSelectorString,
        assertions: assertions.map(assertion => AssertionService.getStringAssertion(assertion)),
        originalSelector: newSelectorString,
        isDraft: true,
        name,
      };

      if (isEditing) {
        CreateAssertionModalAnalyticsService.onCreateAssertionFormSubmit();
        await update(selector, definition);
      } else {
        CreateAssertionModalAnalyticsService.onEditAssertionFormSubmit();
        await add(definition);
      }

      setIsOpen(false);
      onClearMatchedSpans();
    },
    [formProps, onClearMatchedSpans, update, add]
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

export default TestSpecFormProvider;
