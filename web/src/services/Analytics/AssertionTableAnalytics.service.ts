import {Categories} from '../../constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  EditAssertionButtonClick = 'edit-assertion-button-click',
}

const AssertionTableAnalyticsService = () => {
  const onEditAssertionButtonClick = (assertionId: string) => {
    AnalyticsService.event<Actions>(Categories.SpanDetail, Actions.EditAssertionButtonClick, assertionId);
  };

  return {
    onEditAssertionButtonClick,
  };
};

export default AssertionTableAnalyticsService();
