import {MoreOutlined, QuestionCircleOutlined} from '@ant-design/icons';
import {Badge, Dropdown, Menu, Table, Tooltip} from 'antd';
import TestStateBadge from 'components/TestStateBadge';
import {differenceInSeconds} from 'date-fns';
import {FC} from 'react';
import CustomTable from '../../components/CustomTable';
import {Steps} from '../../components/GuidedTour/testDetailsStepList';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';
import {TAssertionResultList} from '../../types/Assertion.types';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {TestState} from '../../constants/TestRunResult.constants';
import {ITestRunResult} from '../../types/TestRunResult.types';

const {onTestRunClick} = TestAnalyticsService;

type TextRowProps = {
  testResultList: ITestRunResult[];
  isLoading: boolean;
  onSelectResult(result: ITestRunResult): void;
};

const validStatusList = [TestState.FINISHED, TestState.AWAITING_TEST_RESULTS];

const getTestResultCount = (assertionResultList: TAssertionResultList, type: 'all' | 'passed' | 'failed' = 'all') => {
  const spanAssertionList = assertionResultList.flatMap(({spanAssertionResults}) => spanAssertionResults);

  if (type === 'all') return spanAssertionList.length;

  return spanAssertionList
    .filter(el => el)
    .filter(({passed}) => {
      switch (type) {
        case 'failed': {
          return !passed;
        }

        case 'passed':
        default: {
          return passed;
        }
      }
    }).length;
};

const TextDetailsTable: FC<TextRowProps> = ({isLoading, onSelectResult, testResultList}) => {
  return (
    <CustomTable
      pagination={{pageSize: 10}}
      rowKey="resultId"
      loading={isLoading}
      dataSource={testResultList?.slice()?.reverse()}
      onRow={record => {
        return {
          onClick: () => {
            const testResult = record as ITestRunResult;
            onTestRunClick(testResult.traceId);
            onSelectResult(testResult);
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
        title={
          <span data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.ExecutionTime)}>
            Execution time
          </span>
        }
        key="executionTime"
        width="10%"
        render={(value, {createdAt, completedAt}: ITestRunResult) => {
          if (!createdAt || !completedAt) return '';
          const executionTime = differenceInSeconds(new Date(completedAt), new Date(createdAt)) + 1;

          return `${executionTime}s`;
        }}
      />
      <Table.Column
        title={<span data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.Status)}>Status</span>}
        key="state"
        width="20%"
        render={(value, {state}: ITestRunResult) => {
          return <TestStateBadge testState={state} />;
        }}
      />
      <Table.Column
        width="5%"
        title={<span data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.Assertions)}>Total</span>}
        key="total"
        dataIndex="state"
        render={(value, {state, assertionResult = []}: ITestRunResult) => {
          if (validStatusList.includes(state)) {
            return getTestResultCount(assertionResult, 'all');
          }

          return '';
        }}
      />
      <Table.Column
        width="5%"
        title={
          <Badge
            data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.Passed)}
            count="P"
            style={{backgroundColor: '#49AA19'}}
          />
        }
        key="passed"
        dataIndex="state"
        render={(value, {state, assertionResult = []}: ITestRunResult) => {
          if (validStatusList.includes(state)) {
            return getTestResultCount(assertionResult, 'passed');
          }

          return '';
        }}
      />
      <Table.Column
        width="5%"
        title={<Badge data-tour={GuidedTourService.getStep(GuidedTours.TestDetails, Steps.Failed)} count="F" />}
        dataIndex="state"
        key="failed"
        render={(value, {state, assertionResult = []}: ITestRunResult) => {
          if (validStatusList.includes(state)) {
            return getTestResultCount(assertionResult, 'failed');
          }

          return '';
        }}
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
