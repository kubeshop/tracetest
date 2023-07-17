import {Dropdown, Menu} from 'antd';
import {useNavigate} from 'react-router-dom';

import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import {ResourceType} from 'types/Resource.type';
import * as S from './TransactionRunActionsMenu.styled';

interface IProps {
  runId: string;
  transactionId: string;
  isRunView?: boolean;
}

const TransactionRunActionsMenu = ({runId, transactionId, isRunView = false}: IProps) => {
  const navigate = useNavigate();
  const onDelete = useDeleteResourceRun({id: transactionId, isRunView, type: ResourceType.Transaction});

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
            <Menu.Item key="automate" onClick={() => navigate(`/transaction/${transactionId}/run/${runId}/automate`)}>
              Automate
            </Menu.Item>
            <Menu.Item
              key="edit"
              onClick={({domEvent}) => {
                domEvent.stopPropagation();
                navigate(`/transaction/${transactionId}/run/${runId}`);
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
          <S.ActionButtonRunView data-cy={`transaction-run-actions-button-${runId}`} />
        ) : (
          <S.ActionButton data-cy={`transaction-run-actions-button-${runId}`} />
        )}
      </Dropdown>
    </span>
  );
};

export default TransactionRunActionsMenu;
