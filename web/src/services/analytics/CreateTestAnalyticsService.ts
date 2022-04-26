import AnalyticsService, { Categories, Labels } from "./AnalyticsService";

enum Actions {
  CreateTestFormSubmit = 'create-test-form-submit',
}

type TCreateTestAnalytics = {
  onCreateTestFormSubmit(): void;
};

const {event} = AnalyticsService(Categories.Home);

const CreateTestAnalyticsService = (): TCreateTestAnalytics => {  
  const onCreateTestFormSubmit = () => {
    event(Actions.CreateTestFormSubmit, Labels.Form);
  };

  return {
    onCreateTestFormSubmit,
  };
};

export default CreateTestAnalyticsService();
