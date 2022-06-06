import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  EditAssertionButtonClick = 'edit-assertion-button-click',
  AssertionClick = 'assertion-click',
  AssertionEditButtonClick = 'assertion-edit-button-click',
  AssertionDeleteButtonClick = 'assertion-delete-button-click',
  AssertionRevert = 'assertion-revert',
}

const AssertionAnalyticsService = () => {
  const onAssertionClick = () => {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.AssertionClick, Labels.Button);
  };

  const onAssertionEdit = () => {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.EditAssertionButtonClick, Labels.Button);
  };

  const onAssertionDelete = () => {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.AssertionDeleteButtonClick, Labels.Button);
  };

  const onRevertAssertion = () => {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.AssertionRevert, Labels.Button);
  };

  return {
    onAssertionClick,
    onAssertionEdit,
    onAssertionDelete,
    onRevertAssertion,
  };
};

export default AssertionAnalyticsService();
