import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useParams} from 'react-router-dom';
import TransactionProvider from '../../providers/Transaction';
import Content from './Content';

const Transaction: React.FC = () => {
  const {transactionId = ''} = useParams();
  return (
    <Layout>
      <TransactionProvider transactionId={transactionId}>
        <Content />
      </TransactionProvider>
    </Layout>
  );
};

export default withAnalytics(Transaction, 'transaction-details');
