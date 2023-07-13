import TransactionRunLayout from 'components/TransactionRunLayout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const TransactionRunTrigger = () => (
  <TransactionRunLayout>
    <Content />
  </TransactionRunLayout>
);

export default withAnalytics(TransactionRunTrigger, 'transaction-details-trigger');
