import {withTracker} from 'ga-4-react';
import Layout from 'components/Layout';
import EditTestModalProvider from '../../components/EditTestModal/EditTestModal.provider';
import TestContent from './TestContent';

const TestPage: React.FC = () => {
  return (
    <Layout>
      <EditTestModalProvider>
        <TestContent />
      </EditTestModalProvider>
    </Layout>
  );
};

export default withTracker(TestPage);
