import {withTracker} from 'ga-4-react';
import Layout from 'components/Layout';
import TestContent from './TestContent';

const TestPage: React.FC = () => {
  return (
    <Layout>
      <TestContent />
    </Layout>
  );
};

export default withTracker(TestPage);
