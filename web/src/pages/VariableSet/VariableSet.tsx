import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import VariableSetContent from './VariableSetContent';

const VariableSet = () => <VariableSetContent />;

export default withAnalytics(VariableSet, 'variable-set');
