import {skipToken} from '@reduxjs/toolkit/dist/query';
import {Table} from 'antd';
import {ColumnsType} from 'antd/lib/table';
import Title from 'antd/lib/typography/Title';
import {useGetTestResultsQuery} from 'services/TestService';
import {ITestResult, TestId} from 'types';

interface IProps {
  testId: TestId;
  onSelectResult: (result: ITestResult) => void;
}

const TestDetails = ({testId, onSelectResult}: IProps) => {
  const {data: testResults, isLoading} = useGetTestResultsQuery(testId ?? skipToken);

  const columns: ColumnsType<ITestResult> = [
    {
      title: 'Test Result Date',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (value, record) => {
        return (
          <p key={record.id}>
            {Intl.DateTimeFormat('default', {dateStyle: 'full', timeStyle: 'medium'} as any).format(new Date(value))}
          </p>
        );
      },
    },
  ];

  return (
    <div style={{overflowY: 'hidden'}}>
      <Title level={5}>Test Results</Title>
      <Table
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
