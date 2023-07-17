import TransactionRunLayout from 'components/TransactionRunLayout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const TransactionRunOverview = () => (
  <TransactionRunLayout>
    <Content />
  </TransactionRunLayout>
);

export default withAnalytics(TransactionRunOverview, 'transaction-details-overview');
