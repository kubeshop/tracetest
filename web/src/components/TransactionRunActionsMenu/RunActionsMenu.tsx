import {Dropdown, Menu} from 'antd';
import useDeleteResourceRun from 'hooks/useDeleteResourceRun';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunActionsMenu.styled';

interface IProps {
  runId: string;
  transactionId: string;
  isRunView?: boolean;
}

const TransactionRunActionsMenu = ({runId, transactionId, isRunView = false}: IProps) => {
  const onDelete = useDeleteResourceRun({id: transactionId, isRunView, type: ResourceType.Transaction});

  return (
    <span className="ant-dropdown-link" onClick={e => e.stopPropagation()} style={{textAlign: 'right'}}>
      <Dropdown
        overlay={
          <Menu>
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
