import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import VariableSetContent from './VariableSetContent';

const VariableSet = () => (
  <Layout hasMenu>
    <VariableSetContent />
  </Layout>
);

export default withAnalytics(VariableSet, 'variable-set');
