import {Dropdown, Menu} from 'antd';
import {useNavigate} from 'react-router-dom';

import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import {ResourceType} from 'types/Resource.type';
import {useFileViewerModal} from '../FileViewerModal/FileViewerModal.provider';
import * as S from './TransactionRunActionsMenu.styled';

interface IProps {
  runId: string;
  transactionId: string;
  isRunView?: boolean;
  transactionVersion: number;
}

const TransactionRunActionsMenu = ({runId, transactionId, isRunView = false, transactionVersion}: IProps) => {
  const {onDefinition} = useFileViewerModal();
  const navigate = useNavigate();
  const onDelete = useDeleteResourceRun({id: transactionId, isRunView, type: ResourceType.Transaction});

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item
              data-cy="view-transaction-definition-button"
              key="view-transaction-definition"
              onClick={() => onDefinition(ResourceType.Transaction, transactionId, transactionVersion)}
            >
              Transaction Definition
            </Menu.Item>
            <Menu.Item
              data-cy="test-edit-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                navigate(`/transaction/${transactionId}/run/${runId}`);
              }}
              key="edit"
            >
              Edit
            </Menu.Item>
            <Menu.Item
              data-cy="test-delete-button"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                onDelete(runId);
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
          <S.ActionButtonRunView data-cy={`transaction-run-actions-button-${runId}`} />
        ) : (
          <S.ActionButton data-cy={`transaction-run-actions-button-${runId}`} />
        )}
      </Dropdown>
    </span>
  );
};

export default TransactionRunActionsMenu;
