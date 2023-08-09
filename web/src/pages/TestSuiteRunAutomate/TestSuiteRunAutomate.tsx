import TestSuiteRunLayout from 'components/TestSuiteRunLayout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const TestSuiteRunAutomate = () => (
  <TestSuiteRunLayout>
    <Content />
  </TestSuiteRunLayout>
);

export default withAnalytics(TestSuiteRunAutomate, 'testsuite-details-automate');
