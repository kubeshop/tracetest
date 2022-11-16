import {useParams} from 'react-router-dom';
import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TransactionRunProvider from 'providers/TransactionRun/TransactionRun.provider';
import TransactionContent from './Content';

const TransactionRunDetail = () => {
  const {transactionId = '', runId = ''} = useParams();

  return (
    <Layout>
      <TransactionRunProvider transactionId={transactionId} runId={runId}>
        <TransactionContent />
      </TransactionRunProvider>
    </Layout>
  );
};

export default withAnalytics(TransactionRunDetail, 'transaction-details');
