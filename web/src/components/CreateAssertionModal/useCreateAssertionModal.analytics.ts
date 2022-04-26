import useAnalytics, {Categories, Labels} from '../Analytics/useAnalytics';

enum Actions {
  EditAssertionFormSubmit = 'edit-assertion-form-submit',
  CreateAssertionFormSubmit = 'create-assertion-form-submit',
  SelectorChange = 'create-assertion-modal-selector-change',
  ChecksChange = 'create-assertion-modal-assertion-checks-change',
  AddCheck = 'create-assertion-modal-add-check',
  RemoveCheck = 'create-assertion-modal-remove-check',
}

type TCreateAssertionModalAnalytics = {
  onEditAssertionFormSubmit(assertionId: string): void;
  onCreateAssertionFormSubmit(testId: string): void;
  onSelectorChange(selector: string): void;
  onChecksChange(checks: string): void;
  onAddCheck(): void;
  onRemoveCheck(): void;
};

const useCreateAssertionModalAnalytics = (): TCreateAssertionModalAnalytics => {
  const {event} = useAnalytics(Categories.Assertion);

  const onCreateAssertionFormSubmit = (testId: string) => {
    event(Actions.CreateAssertionFormSubmit, testId);
  };

  const onEditAssertionFormSubmit = (assertionId: string) => {
    event(Actions.EditAssertionFormSubmit, assertionId);
  };

  const onSelectorChange = (selector: string) => {
    event(Actions.SelectorChange, selector);
  };

  const onChecksChange = (checks: string) => {
    event(Actions.SelectorChange, checks);
  };

  const onAddCheck = () => {
    event(Actions.AddCheck, Labels.Button);
  };

  const onRemoveCheck = () => {
    event(Actions.RemoveCheck, Labels.Button);
  };

  return {
    onCreateAssertionFormSubmit,
    onEditAssertionFormSubmit,
    onSelectorChange,
    onChecksChange,
    onAddCheck,
    onRemoveCheck,
  };
};

export default useCreateAssertionModalAnalytics;
