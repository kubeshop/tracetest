import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import {useParams} from 'react-router-dom';
import TransactionRunDetailProvider from '../../providers/TransactionRunDetail';
import TransactionContent from './Content';

const TransactionRunDetail: React.FC = () => {
  const {transactionId = '', runId = ''} = useParams();
  return (
    <Layout>
      <TransactionRunDetailProvider transactionId={transactionId} runId={runId}>
        <TransactionContent />
      </TransactionRunDetailProvider>
    </Layout>
  );
};

export default withAnalytics(TransactionRunDetail, 'transaction-details');
