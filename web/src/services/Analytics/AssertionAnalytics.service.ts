import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  EditAssertionButtonClick = 'edit-assertion-button-click',
  AssertionClick = 'assertion-click',
  AssertionEditButtonClick = 'assertion-edit-button-click',
  AssertionDeleteButtonClick = 'assertion-delete-button-click',
  AssertionRevert = 'assertion-revert',
}

const AssertionAnalyticsService = () => ({
  onAssertionClick() {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.AssertionClick, Labels.Button);
  },
  onAssertionEdit() {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.EditAssertionButtonClick, Labels.Button);
  },
  onAssertionDelete() {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.AssertionDeleteButtonClick, Labels.Button);
  },
  onRevertAssertion() {
    AnalyticsService.event<Actions>(Categories.TestResults, Actions.AssertionRevert, Labels.Button);
  },
});

export default AssertionAnalyticsService();
