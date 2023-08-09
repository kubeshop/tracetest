import TestSuiteRunLayout from 'components/TestSuiteRunLayout';
import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Content from './Content';

const TestSuiteRunOverview = () => (
  <TestSuiteRunLayout>
    <Content />
  </TestSuiteRunLayout>
);

export default withAnalytics(TestSuiteRunOverview, 'testsuite-run-details-overview');
