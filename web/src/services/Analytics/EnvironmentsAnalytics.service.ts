import {Categories, Labels} from 'constants/Analytics.constants';
import AnalyticsService from './Analytics.service';

export enum Actions {
  CreateEnvironmentClick = 'create-environment-button-click',
  EnvironmentClick = 'environment-click',
}

type TEnvironmentsAnalytics = {
  onCreateEnvironmentClick(): void;
  onEnvironmentClick(environmentId: string): void;
};

const HomeAnalyticsService = (): TEnvironmentsAnalytics => {
  return {
    onCreateEnvironmentClick: () => {
      AnalyticsService.event(Categories.Environments, Actions.CreateEnvironmentClick, Labels.Button);
    },
    onEnvironmentClick: (environmentId: string) => {
      AnalyticsService.event(Categories.Environments, Actions.EnvironmentClick, environmentId);
    },
  };
};

export default HomeAnalyticsService();
