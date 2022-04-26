import useAnalytics, {Categories, Labels} from '../Analytics/useAnalytics';

enum Actions {
  CreateTestFormSubmit = 'create-test-form-submit',
}

type TCreateTestAnalytics = {
  onCreateTestFormSubmit(): void;
};

const useCreateTestAnalytics = (): TCreateTestAnalytics => {
  const {event} = useAnalytics(Categories.Home);

  const onCreateTestFormSubmit = () => {
    event(Actions.CreateTestFormSubmit, Labels.Form);
  };

  return {
    onCreateTestFormSubmit,
  };
};

export default useCreateTestAnalytics;
