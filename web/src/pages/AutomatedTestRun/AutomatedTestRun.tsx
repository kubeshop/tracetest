import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import TestSpecFormProvider from 'components/TestSpecForm/TestSpecForm.provider';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test/Test.provider';
import Content from './Content';

const AutomatedTestRun = () => {
  const {testId = '', version = '1'} = useParams();

  return (
    <Layout hasMenu>
      <TestProvider testId={testId} version={Number(version)}>
        <TestSpecFormProvider testId={testId}>
          <Content />
        </TestSpecFormProvider>
      </TestProvider>
    </Layout>
  );
};

export default withAnalytics(AutomatedTestRun, 'automated-test-run');
