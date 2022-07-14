import Layout from 'components/Layout';
import {withTracker} from 'ga-4-react';
import {useParams} from 'react-router-dom';
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

export default withTracker(EditTestPage);
