import {useParams, useSearchParams} from 'react-router-dom';

import TestSpecFormProvider from 'components/TestSpecForm/TestSpecForm.provider';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test/Test.provider';
import Content from './Content';

const AutomatedTestRun = () => {
  const {testId = ''} = useParams();
  const [query] = useSearchParams();
  const version = query.get('version') ? Number(query.get('version')) : undefined;

  return (
    <TestProvider testId={testId} version={version}>
      <TestSpecFormProvider testId={testId}>
        <Content />
      </TestSpecFormProvider>
    </TestProvider>
  );
};

export default withAnalytics(AutomatedTestRun, 'automated-test-run');
