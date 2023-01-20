import {Dropdown, Menu} from 'antd';
import {useCallback} from 'react';

import * as S from './ResourceCard.styled';

interface IProps {
  id: string;
  canEdit: boolean;
  onDelete(): void;
  onEdit(): void;
}

const ResourceCardActions = ({id, canEdit = true, onDelete, onEdit}: IProps) => {
  const onDeleteClick = useCallback(
    ({domEvent}) => {
      domEvent?.stopPropagation();
      onDelete();
    },
    [onDelete]
  );

  const onEditClick = useCallback(
    ({domEvent}) => {
      domEvent?.stopPropagation();
      onEdit();
    },
    [onEdit]
  );

  const menuItems = [];
  if (canEdit) {
    menuItems.push({key: 'edit', label: <span data-cy="test-card-edit">Edit</span>, onClick: onEditClick});
  }
  menuItems.push({key: 'delete', label: <span data-cy="test-card-delete">Delete</span>, onClick: onDeleteClick});

  return (
    <Dropdown
      overlay={<Menu items={menuItems} />}
      placement="bottomLeft"
      trigger={['click']}
    >
      <span data-cy={`test-actions-button-${id}`} className="ant-dropdown-link" onClick={e => e.stopPropagation()}>
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default ResourceCardActions;
