import {Dropdown, Menu} from 'antd';

import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {ResourceType} from 'types/Resource.type';
import * as S from './TestSuiteRunActionsMenu.styled';

interface IProps {
  runId: string;
  testSuiteId: string;
  isRunView?: boolean;
}

const TestSuiteRunActionsMenu = ({runId, testSuiteId, isRunView = false}: IProps) => {
  const {navigate} = useDashboard();
  const onDelete = useDeleteResourceRun({id: testSuiteId, isRunView, type: ResourceType.TestSuite});

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item key="automate" onClick={() => navigate(`/testsuite/${testSuiteId}/run/${runId}/automate`)}>
              Automate
            </Menu.Item>
            <Menu.Item
              key="edit"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                navigate(`/testsuite/${testSuiteId}/run/${runId}`);
              }}
            >
              Edit
            </Menu.Item>
            <Menu.Item
              key="delete"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                onDelete(runId);
              }}
            >
              Delete
            </Menu.Item>
          </Menu>
        }
        placement="bottomLeft"
        trigger={['click']}
      >
        {isRunView ? (
          <S.ActionButtonRunView data-cy={`testsuite-run-actions-button-${runId}`} />
        ) : (
          <S.ActionButton data-cy={`testsuite-run-actions-button-${runId}`} />
        )}
      </Dropdown>
    </span>
  );
};

export default TestSuiteRunActionsMenu;
