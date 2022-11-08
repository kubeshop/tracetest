import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TransactionProvider from 'providers/Transaction';
import TransactionRunProvider from 'providers/TransactionRun';
import {useParams} from 'react-router-dom';
import TransactionContent from './Content';

const TransactionRunDetail = () => {
  const {transactionId = '', runId = ''} = useParams();
  return (
    <Layout>
      <TransactionRunProvider transactionId={transactionId} runId={runId}>
        <TransactionProvider transactionId={transactionId}>
          <TransactionContent />
        </TransactionProvider>
      </TransactionRunProvider>
    </Layout>
  );
};

export default withAnalytics(TransactionRunDetail, 'transaction-details');
