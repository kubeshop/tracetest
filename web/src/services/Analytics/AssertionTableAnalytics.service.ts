import AnalyticsService, { Categories } from './Analytics.service';

enum Actions {
  EditAssertionButtonClick = 'edit-assertion-button-click',
}

type TAssertionTableAnalytics = {
  onEditAssertionButtonClick(assertionId: string): void;
};

const {event} = AnalyticsService(Categories.SpanDetail);

const AssertionTableAnalyticsService = (): TAssertionTableAnalytics => {
  const onEditAssertionButtonClick = (assertionId: string) => {
    event(Actions.EditAssertionButtonClick, assertionId);
  };

  return {
    onEditAssertionButtonClick,
  };
};

export default AssertionTableAnalyticsService();
