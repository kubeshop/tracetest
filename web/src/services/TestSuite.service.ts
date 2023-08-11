import {TDraftTestSuite} from 'types/TestSuite.types';
import TestSuite, {TRawTestSuiteResource} from 'models/TestSuite.model';

const TestSuiteService = () => ({
  getRawFromDraft(draft: TDraftTestSuite): TRawTestSuiteResource {
    return {
      spec: {...draft, fullSteps: draft.steps?.map(step => ({id: step}))},
      type: 'TestSuite',
    };
  },

  getInitialValues(suite: TestSuite): TDraftTestSuite {
    return {
      ...suite,
      steps: suite.fullSteps.map(step => step.id),
    };
  },
});

export default TestSuiteService();
