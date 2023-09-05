import {Dropdown, Menu} from 'antd';
import {useCallback, useMemo} from 'react';

import {Operation, useCustomization} from 'providers/Customization';
import * as S from './ResourceCard.styled';

interface IProps {
  id: string;
  shouldEdit: boolean;
  onDelete(): void;
  onEdit(): void;
}

const ResourceCardActions = ({id, shouldEdit = true, onDelete, onEdit}: IProps) => {
  const {getIsAllowed} = useCustomization();

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

  const menuItems = useMemo(() => {
    const defaultItems = [
      {
        key: 'delete',
        label: <span data-cy="test-card-delete">Delete</span>,
        onClick: onDeleteClick,
        disabled: !getIsAllowed(Operation.Edit),
      },
    ];

    return shouldEdit
      ? [
          {
            key: 'edit',
            label: <span data-cy="test-card-edit">Edit</span>,
            onClick: onEditClick,
            disabled: !getIsAllowed(Operation.Edit),
          },
          ...defaultItems,
        ]
      : defaultItems;
  }, [getIsAllowed, onDeleteClick, onEditClick, shouldEdit]);

  return (
    <Dropdown overlay={<Menu items={menuItems} />} placement="bottomLeft" trigger={['click']}>
      <span data-cy={`test-actions-button-${id}`} className="ant-dropdown-link" onClick={e => e.stopPropagation()}>
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default ResourceCardActions;
