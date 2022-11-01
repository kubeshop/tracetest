import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import EnvironmentContent from './EnvironmentContent';

const Environment = () => (
  <Layout hasMenu>
    <EnvironmentContent />
  </Layout>
);

export default withAnalytics(Environment, 'environments');
