import {Dropdown, Menu} from 'antd';
import {Link} from 'react-router-dom';
import useDeleteTestRun from '../../hooks/useDeleteTestRun';
import {useFileViewerModal} from '../FileViewerModal/FileViewerModal.provider';
import * as S from './RunActionsMenu.styled';

interface IProps {
  resultId: string;
  testId: string;
  testVersion: number;
  isRunView?: boolean;
}

const RunActionsMenu = ({resultId, testId, testVersion, isRunView = false}: IProps) => {
  const {loadJUnit, loadTestDefinitionYaml} = useFileViewerModal();

  const onDelete = useDeleteTestRun({isRunView, testId});

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
            {testId && (
              <Menu.Item data-cy="test-edit-button" key="edit">
                <Link to={`/test/${testId}/edit`}>Edit Test</Link>
              </Menu.Item>
            )}
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
        <S.ActionButton data-cy={`result-actions-button-${resultId}`} />
      </Dropdown>
    </span>
  );
};

export default RunActionsMenu;
