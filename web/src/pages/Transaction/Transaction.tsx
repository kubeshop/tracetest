import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useParams} from 'react-router-dom';
import TransactionProvider from '../../providers/TransactionRunDetail';
import TransactionContent from './Content';

const Transaction: React.FC = () => {
  const {transactionId = '', runId = ''} = useParams();
  return (
    <Layout>
      <TransactionProvider transactionId={transactionId} runId={runId}>
        <TransactionContent />
      </TransactionProvider>
    </Layout>
  );
};

export default withAnalytics(Transaction, 'transaction-details');
