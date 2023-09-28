import {Dropdown, Menu} from 'antd';

import {useFileViewerModal} from 'components/FileViewerModal/FileViewerModal.provider';
import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import {Operation, useCustomization} from 'providers/Customization';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunActionsMenu.styled';

interface IProps {
  resultId: number;
  testId: string;
  isRunView?: boolean;
  testSuiteId?: string;
  testSuiteRunId: number;
  origin?: string;
}

const RunActionsMenu = ({resultId, testId, testSuiteId, testSuiteRunId, isRunView = false, origin}: IProps) => {
  const {getIsAllowed} = useCustomization();
  const {onJUnit} = useFileViewerModal();
  const {navigate} = useDashboard();
  const onDelete = useDeleteResourceRun({id: testId, isRunView, type: ResourceType.Test});

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
            {!!testSuiteId && !!testSuiteRunId && (
              <Menu.Item
                data-cy="testsuite-run-button"
                key="testsuite-run"
                onClick={() => {
                  navigate(`/testsuite/${testSuiteId}/run/${testSuiteRunId}`);
                }}
              >
                Test Suite Run
              </Menu.Item>
            )}
            <Menu.Item
              data-cy="view-junit-button"
              key="view-junit"
              onClick={() => {
                TestRunAnalyticsService.onLoadJUnitReport();
                onJUnit(testId, resultId);
              }}
            >
              JUnit Results
            </Menu.Item>
            <Menu.Item
              data-cy="automate-test-button"
              key="automate-test"
              onClick={() => {
                navigate(`/test/${testId}/run/${resultId}/automate`, {state: {origin}});
              }}
            >
              Automate
            </Menu.Item>
            <Menu.Item
              data-cy="test-edit-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                navigate(`/test/${testId}/run/${resultId}`, {state: {origin}});
              }}
              key="edit"
              disabled={!getIsAllowed(Operation.Edit)}
            >
              Edit
            </Menu.Item>
            <Menu.Item
              data-cy="test-delete-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                onDelete(resultId);
              }}
              key="delete"
              disabled={!getIsAllowed(Operation.Edit)}
            >
              Delete
            </Menu.Item>
          </Menu>
        }
        placement="bottomLeft"
        trigger={['click']}
      >
        {isRunView ? (
          <S.ActionButtonRunView data-cy={`result-actions-button-${resultId}`} />
        ) : (
          <S.ActionButton data-cy={`result-actions-button-${resultId}`} />
        )}
      </Dropdown>
    </span>
  );
};

export default RunActionsMenu;
