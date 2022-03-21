import {skipToken} from '@reduxjs/toolkit/dist/query';
import {Button, Table} from 'antd';
import {ColumnsType} from 'antd/lib/table';
import Title from 'antd/lib/typography/Title';
import {useGetTestResultsQuery, useRunTestMutation} from 'services/TestService';
import {ITestResult, TestId} from 'types';

interface IProps {
  testId: TestId;
  onSelectResult: (result: ITestResult) => void;
}

const TestDetails = ({testId, onSelectResult}: IProps) => {
  const {data: testResults, isLoading} = useGetTestResultsQuery(testId ?? skipToken);
  const [runTest, result] = useRunTestMutation();
  const columns: ColumnsType<ITestResult> = [
    {
      title: 'Test Result Date',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (value, record) => {
        return (
          <p key={record.resultId}>
            {Intl.DateTimeFormat('default', {dateStyle: 'full', timeStyle: 'medium'} as any).format(new Date(value))}
          </p>
        );
      },
    },
    {
      title: 'Deployment Version',
      dataIndex: 'id',
      key: 'id',
      render: (value, record) => {
        return <p key={record.id}>0.0.1</p>;
      },
    },
    {
      title: 'Environment',
      dataIndex: 'id',
      key: 'id',
      render: (value, record) => {
        return <p key={record.id}>Staging</p>;
      },
    },
    {
      title: 'Assertion Result',
      dataIndex: 'id',
      key: 'id',
      render: (value, record) => {
        return <p key={record.id}>Passed</p>;
      },
    },
  ];

  const handleRunTest = () => {
    if (testId) {
      runTest(testId);
    }
  };

  return (
    <div style={{overflowY: 'hidden'}}>
      <Button style={{marginBottom: 24, marginTop: 8}} shape="round" onClick={handleRunTest} loading={result.isLoading}>
        Generate Trace
      </Button>

      <Title style={{marginBottom: 8}} level={5}>
        Test Results
      </Title>
      <Table
        pagination={{pageSize: 5}}
        rowKey="id"
        loading={isLoading}
        columns={columns}
        dataSource={testResults?.slice()?.reverse()}
        onRow={record => {
          return {
            onClick: () => {
              onSelectResult(record);
            },
          };
        }}
      />
    </div>
  );
};

export default TestDetails;
