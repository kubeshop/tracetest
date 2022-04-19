import {MoreOutlined, QuestionCircleOutlined} from '@ant-design/icons';
import {Badge, Dropdown, Menu, Table, Tooltip} from 'antd';
import {differenceInSeconds} from 'date-fns';
import {FC} from 'react';
import CustomTable from '../../components/CustomTable';
import {Assertion, ITestResult, TestState} from '../../types';

const BadgeStatusMap: Record<
  TestState,
  {status: 'processing' | 'success' | 'error' | 'default' | 'warning' | undefined; label: string}
> = {
  [TestState.CREATED]: {
    status: 'default',
    label: 'Created',
  },
  [TestState.EXECUTING]: {
    status: 'processing',
    label: 'Running',
  },
  [TestState.AWAITING_TRACE]: {
    status: 'warning',
    label: 'Awaiting trace',
  },
  [TestState.AWAITING_TEST_RESULTS]: {
    status: 'success',
    label: 'Awaiting test results',
  },
  [TestState.FINISHED]: {
    status: 'success',
    label: 'Finished',
  },
  [TestState.FAILED]: {
    status: 'error',
    label: 'Failed executing test run',
  },
};

type TextRowProps = {
  assertionList: Assertion[];
  testResultList: ITestResult[];
  isLoading: boolean;
  onSelectResult(result: ITestResult): void;
};

const validStatusList = [TestState.FINISHED, TestState.AWAITING_TEST_RESULTS];

const TextDetailsTable: FC<TextRowProps> = ({assertionList, isLoading, onSelectResult, testResultList}) => {
  return (
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
        title="Time"
        dataIndex="createdAt"
        key="createdAt"
        width="30%"
        render={value =>
          Intl.DateTimeFormat('default', {dateStyle: 'full', timeStyle: 'medium'} as any).format(new Date(value))
        }
      />
      <Table.Column
        title="Execution time"
        key="executionTime"
        width="10%"
        render={(value, {createdAt, completedAt}: ITestResult) => {
          if (!createdAt || !completedAt) return '';
          const executionTime = differenceInSeconds(new Date(createdAt), new Date(completedAt)) + 1;

          return `${executionTime}s`;
        }}
      />
      <Table.Column
        title="Status"
        key="status"
        width="20%"
        render={(value, {state}: ITestResult) => {
          const {status, label} = BadgeStatusMap[state] || BadgeStatusMap.CREATED;

          return <Badge status={status} text={label} />;
        }}
      />
      <Table.Column
        width="5%"
        title="Total"
        key="total"
        dataIndex="state"
        render={state => (validStatusList.includes(state) ? assertionList?.length ?? 0 : '')}
      />
      <Table.Column
        width="3%"
        title={<Badge count="P" style={{backgroundColor: '#49AA19'}} />}
        key="passed"
        dataIndex="state"
        render={state => (validStatusList.includes(state) ? 0 : '')}
      />
      <Table.Column
        width="3%"
        title={<Badge count="F" />}
        dataIndex="state"
        key="failed"
        render={state => (validStatusList.includes(state) ? 0 : '')}
      />
      <Table.Column
        width="3%"
        title={
          <Tooltip title="The number of Total/Pass/Fail assertions">
            <QuestionCircleOutlined style={{color: '#8C8C8C'}} />
          </Tooltip>
        }
        key="question"
      />
      <Table.Column
        title="Actions"
        key="actions"
        align="right"
        render={() => {
          const menuLayout = (
            <Menu>
              <Menu.Item key="delete">Delete</Menu.Item>
            </Menu>
          );

          return (
            <Dropdown overlay={menuLayout} placement="bottomLeft" trigger={['click']}>
              <span
                className="ant-dropdown-link"
                onClick={e => {
                  e.preventDefault();
                  e.stopPropagation();
                }}
              >
                <MoreOutlined style={{fontSize: 24}} />
              </span>
            </Dropdown>
          );
        }}
      />
    </CustomTable>
  );
};

export default TextDetailsTable;
