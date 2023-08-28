import {useParams} from 'react-router-dom';

import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import TestProvider from 'providers/Test';
import Content from './Content';

const Test = () => {
  const {testId = ''} = useParams();

  return (
    <TestProvider testId={testId}>
      <Content />
    </TestProvider>
  );
};

export default withAnalytics(Test, 'test-details');
