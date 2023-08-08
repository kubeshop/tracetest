import {useParams} from 'react-router-dom';
import Layout from 'components/Layout';
import TransactionHeader from 'components/TransactionHeader';
import TransactionRunProvider from 'providers/TransactionRun/TransactionRun.provider';

interface IProps {
  children: React.ReactNode;
}

const TransactionRunLayout = ({children}: IProps) => {
  const {transactionId = '', runId = ''} = useParams();

  return (
    <Layout>
      <TransactionRunProvider transactionId={transactionId} runId={runId}>
        <TransactionHeader />
        {children}
      </TransactionRunProvider>
    </Layout>
  );
};

export default TransactionRunLayout;
