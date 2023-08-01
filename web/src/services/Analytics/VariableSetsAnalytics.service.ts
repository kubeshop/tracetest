import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  CreateVariableSetClick = 'create-variableSet-button-click',
  VariableSetClick = 'variableSet-click',
}

type TVariableSetsAnalytics = {
  onCreateVariableSetClick(): void;
  onVariableSetClick(VariableSetId: string): void;
};

const VariableSetsAnalyticsService = (): TVariableSetsAnalytics => {
  return {
    onCreateVariableSetClick: () => {
      AnalyticsService.event(Categories.VariableSet, Actions.CreateVariableSetClick, Labels.Button);
    },
    onVariableSetClick: (VariableSetId: string) => {
      AnalyticsService.event(Categories.VariableSet, Actions.VariableSetClick, VariableSetId);
    },
  };
};

export default VariableSetsAnalyticsService();
