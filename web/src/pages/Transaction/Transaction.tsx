import {useParams} from 'react-router-dom';

import Layout from 'components/Layout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TransactionProvider from 'providers/Transaction';
import Content from './Content';

const Transaction = () => {
  const {transactionId = ''} = useParams();

  return (
    <Layout hasMenu>
      <TransactionProvider transactionId={transactionId}>
        <Content />
      </TransactionProvider>
    </Layout>
  );
};

export default withAnalytics(Transaction, 'transaction-details');
