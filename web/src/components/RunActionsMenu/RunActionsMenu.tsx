import {Dropdown, Menu} from 'antd';
import {useCallback} from 'react';
import {useNavigate} from 'react-router-dom';
import {useDeleteRunByIdMutation} from 'redux/apis/TraceTest.api';
import TestAnalyticsService from 'services/Analytics/TestAnalytics.service';
import {TTest} from 'types/Test.types';
import {useEditTestModal} from '../EditTestModal/EditTestModal.provider';
import {useFileViewerModal} from '../FileViewerModal/FileViewerModal.provider';
import * as S from './RunActionsMenu.styled';

interface IProps {
  resultId: string;
  testId: string;
  testVersion: number;
  test?: TTest;
  isRunView?: boolean;
}

const RunActionsMenu = ({resultId, testId, testVersion, test, isRunView = false}: IProps) => {
  const {loadJUnit, loadTestDefinitionYaml} = useFileViewerModal();
  const [deleteRunById] = useDeleteRunByIdMutation();
  const navigate = useNavigate();
  const {open} = useEditTestModal();

  const handleOnDelete = useCallback(() => {
    TestAnalyticsService.onDeleteTestRun();
    deleteRunById({testId, runId: resultId});
    if (isRunView) navigate(`/test/${testId}`);
  }, [deleteRunById, isRunView, navigate, resultId, testId]);

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item data-cy="view-junit-button" key="view-junit" onClick={() => loadJUnit(testId, resultId)}>
              JUnit Results
            </Menu.Item>
            <Menu.Item
              data-cy="view-test-definition-button"
              key="view-test-definition"
              onClick={() => loadTestDefinitionYaml(testId, testVersion)}
            >
              Test Definition
            </Menu.Item>
            {test && (
              <Menu.Item
                data-cy="test-edit-button"
                onClick={({domEvent}) => {
                  domEvent.stopPropagation();
                  open(test);
                }}
                key="edit"
              >
                Edit Test
              </Menu.Item>
            )}
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
        <S.ActionButton data-cy={`result-actions-button-${resultId}`} />
      </Dropdown>
    </span>
  );
};

export default RunActionsMenu;
