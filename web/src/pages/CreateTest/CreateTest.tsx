import Layout from 'components/Layout';
import {withTracker} from 'ga-4-react';
import CreateTestProvider from 'providers/CreateTest/CreateTest.provider';
import CreateTestContent from './CreateTestContent';

const CreateTestPage = () => {
  return (
    <Layout>
      <CreateTestProvider>
        <CreateTestContent />
      </CreateTestProvider>
    </Layout>
  );
};

export default withTracker(CreateTestPage);
