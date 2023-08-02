import {useParams, useSearchParams} from 'react-router-dom';

import Layout from 'components/Layout';
import TestSpecFormProvider from 'components/TestSpecForm/TestSpecForm.provider';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test/Test.provider';
import Content from './Content';

const AutomatedTestRun = () => {
  const {testId = ''} = useParams();
  const [query] = useSearchParams();
  const version = query.get('version') ? Number(query.get('version')) : undefined;

  return (
    <Layout hasMenu>
      <TestProvider testId={testId} version={version}>
        <TestSpecFormProvider>
          <Content />
        </TestSpecFormProvider>
      </TestProvider>
    </Layout>
  );
};

export default withAnalytics(AutomatedTestRun, 'automated-test-run');
