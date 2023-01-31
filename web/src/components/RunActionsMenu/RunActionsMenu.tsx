import {Dropdown, Menu} from 'antd';
import {useNavigate} from 'react-router-dom';

import {useFileViewerModal} from 'components/FileViewerModal/FileViewerModal.provider';
import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunActionsMenu.styled';

interface IProps {
  resultId: string;
  testId: string;
  testVersion: number;
  isRunView?: boolean;
  transactionId?: string;
  transactionRunId: string;
}

const RunActionsMenu = ({
  resultId,
  testId,
  testVersion,
  transactionId,
  transactionRunId,
  isRunView = false,
}: IProps) => {
  const {loadJUnit, loadDefinition} = useFileViewerModal();

  const navigate = useNavigate();

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
                loadJUnit(testId, resultId);
              }}
            >
              JUnit Results
            </Menu.Item>
            <Menu.Item
              data-cy="view-test-definition-button"
              key="view-test-definition"
              onClick={() => {
                TestRunAnalyticsService.onLoadTestDefinition();
                loadDefinition(ResourceType.Test, testId, testVersion);
              }}
            >
              Test Definition
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
