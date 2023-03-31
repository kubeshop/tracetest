import TestRunEvent, {TestRunStage} from 'models/TestRunEvent.model';

const TestRunService = () => ({
  getTestRunEventsByStage(events: TestRunEvent[], stage: TestRunStage) {
    return events.filter(event => event.stage === stage);
  },
});

export default TestRunService();
