import useAnalytics, {Categories} from '../Analytics/useAnalytics';

enum Actions {
  EditAssertionButtonClick = 'edit-assertion-button-click',
}

type TAssertionTableAnalytics = {
  onEditAssertionButtonClick(assertionId: string): void;
};

const useAssertionTableAnalytics = (): TAssertionTableAnalytics => {
  const {event} = useAnalytics(Categories.SpanDetail);

  const onEditAssertionButtonClick = (assertionId: string) => {
    event(Actions.EditAssertionButtonClick, assertionId);
  };

  return {
    onEditAssertionButtonClick,
  };
};

export default useAssertionTableAnalytics;
