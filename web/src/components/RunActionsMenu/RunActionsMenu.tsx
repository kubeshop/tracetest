import {Dropdown, Menu} from 'antd';

import {useFileViewerModal} from 'components/FileViewerModal/FileViewerModal.provider';
import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunActionsMenu.styled';

interface IProps {
  resultId: string;
  testId: string;
  isRunView?: boolean;
  transactionId?: string;
  transactionRunId: string;
}

const RunActionsMenu = ({resultId, testId, transactionId, transactionRunId, isRunView = false}: IProps) => {
  const {onJUnit} = useFileViewerModal();
  const {navigate} = useDashboard();
  const onDelete = useDeleteResourceRun({id: testId, isRunView, type: ResourceType.Test});

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
            {!!transactionId && !!transactionRunId && (
              <Menu.Item
                data-cy="transaction-run-button"
                key="transaction-run"
                onClick={() => {
                  navigate(`/transaction/${transactionId}/run/${transactionRunId}`);
                }}
              >
                Transaction Run
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
                navigate(`/test/${testId}/run/${resultId}/automate`);
              }}
            >
              Automate
            </Menu.Item>
            <Menu.Item
              data-cy="test-edit-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                navigate(`/test/${testId}/run/${resultId}`);
              }}
              key="edit"
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
