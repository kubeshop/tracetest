import TransactionRunLayout from 'components/TransactionRunLayout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const TransactionRunAutomate = () => (
  <TransactionRunLayout>
    <Content />
  </TransactionRunLayout>
);

export default withAnalytics(TransactionRunAutomate, 'transaction-details-automate');
