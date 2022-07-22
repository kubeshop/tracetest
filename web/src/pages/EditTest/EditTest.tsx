import Layout from 'components/Layout';
import {useParams} from 'react-router-dom';
import withAnalytics from '../../components/WithAnalytics/WithAnalytics';
import {useGetTestByIdQuery} from '../../redux/apis/TraceTest.api';
import EditTestContent from './EditTestContent';

const EditTestPage = () => {
  const {testId = ''} = useParams();
  const {data: test} = useGetTestByIdQuery({testId});

  return test ? (
    <Layout>
      <EditTestContent test={test} />
    </Layout>
  ) : null;
};

export default withAnalytics(EditTestPage, 'edit-test');
