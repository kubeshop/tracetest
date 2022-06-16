import {Dropdown, Menu} from 'antd';
import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useDeleteRunByIdMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {useFileViewerModal} from '../FileViewerModal/FileViewerModal.provider';
import * as S from './RunActionsDropdown.styled';

interface IProps {
  resultId: string;
  testId: string;
  isRunView?: boolean;
}

const RunActionsDropdown = ({resultId, testId, isRunView = false}: IProps) => {
  const {loadJUnit, loadTestDefinitionYaml} = useFileViewerModal();
  const [deleteRunById] = useDeleteRunByIdMutation();
  const navigate = useNavigate();

  const handleOnDelete = useCallback(() => {
    TestAnalyticsService.onDeleteTestRun();
    deleteRunById({testId, runId: resultId});
    if (isRunView) navigate(`/test/${testId}`);
  }, [deleteRunById, isRunView, navigate, resultId, testId]);

  return (
    <span
      data-cy={`result-actions-button-${resultId}`}
      className="ant-dropdown-link"
      onClick={e => e.stopPropagation()}
      style={{textAlign: 'right'}}
    >
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item data-cy="view-junit-button" key="view-junit" onClick={() => loadJUnit(testId, resultId)}>
              JUnit Results
            </Menu.Item>
            <Menu.Item
              data-cy="view-test-definition-button"
              key="view-test-definition"
              onClick={() => loadTestDefinitionYaml(testId, 1)}
            >
              Test Definition
            </Menu.Item>
            <Menu.Item
              data-cy="test-delete-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                handleOnDelete();
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
        <S.ActionButton />
      </Dropdown>
    </span>
  );
};

export default RunActionsDropdown;
