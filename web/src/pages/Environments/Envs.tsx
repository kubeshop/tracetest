import Layout from 'components/Layout';
import EnvironmentContent from './EnvironmentContent';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';

const Envs = (): JSX.Element => {
  return (
    <Layout hasMenu>
      <EnvironmentContent />
    </Layout>
  );
};

export default withAnalytics(Envs, 'environments');
