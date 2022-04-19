import {Button, Table, Typography} from 'antd';
import {FC, useCallback} from 'react';
import {useRunTestMutation} from 'services/TestService';
import {ITestResult, TestId} from 'types';
import CustomTable from '../../components/CustomTable';
import * as S from './Test.styled';

type TTestDetailsProps = {
  testId: TestId;
  url?: string;
  onSelectResult: (result: ITestResult) => void;
  testResultList: ITestResult[];
  isLoading: boolean;
};

const TestDetails: FC<TTestDetailsProps> = ({testId, testResultList, isLoading, onSelectResult, url}) => {
  const [runTest, result] = useRunTestMutation();

  const handleRunTest = useCallback(() => {
    if (testId) runTest(testId);
  }, [runTest, testId]);

  return (
    <>
      <S.TestDetailsHeader>
        <Typography.Title level={5}>{url}</Typography.Title>
        <Button onClick={handleRunTest} loading={result.isLoading} type="primary" ghost>
          Run Test
        </Button>
      </S.TestDetailsHeader>
      <CustomTable
        pagination={{pageSize: 10}}
        rowKey="resultId"
        loading={isLoading}
        dataSource={testResultList?.slice()?.reverse()}
        onRow={record => {
          return {
            onClick: () => {
              onSelectResult(record as ITestResult);
            },
          };
        }}
      >
        <Table.Column
          title="Test Results"
          dataIndex="createdAt"
          key="createdAt"
          width="30%"
          render={value =>
            Intl.DateTimeFormat('default', {dateStyle: 'full', timeStyle: 'medium'} as any).format(new Date(value))
          }
        />
        <Table.Column title="Assertion Result" dataIndex="url" key="url" width="70%" render={() => 'Passed'} />
      </CustomTable>
    </>
  );
};

export default TestDetails;
