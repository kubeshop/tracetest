import {Dropdown, Menu} from 'antd';
import {useCallback, useMemo} from 'react';

import {Operation, useCustomization} from 'providers/Customization';
import * as S from './ResourceCard.styled';

interface IProps {
  id: string;
  shouldEdit: boolean;
  onDelete(): void;
  onEdit(): void;
  onDuplicate(): void;
}

const ResourceCardActions = ({id, shouldEdit = true, onDelete, onEdit, onDuplicate}: IProps) => {
  const {getIsAllowed} = useCustomization();
  const canEdit = getIsAllowed(Operation.Edit);

  const onAction = useCallback(
    action =>
      ({domEvent}: {domEvent: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>}) => {
        domEvent?.stopPropagation();
        action();
      },
    []
  );

  const menuItems = useMemo(() => {
    const defaultItems = [
      {
        key: 'duplicate',
        label: <span data-cy="test-card-duplicate">Duplicate</span>,
        onClick: onAction(onDuplicate),
        disabled: !canEdit,
      },
      {
        key: 'delete',
        label: <span data-cy="test-card-delete">Delete</span>,
        onClick: onAction(onDelete),
        disabled: !getIsAllowed(Operation.Edit),
      },
    ];

    return shouldEdit
      ? [
          {
            key: 'edit',
            label: <span data-cy="test-card-edit">Edit</span>,
            onClick: onAction(onEdit),
            disabled: !getIsAllowed(Operation.Edit),
          },
          ...defaultItems,
        ]
      : defaultItems;
  }, [canEdit, getIsAllowed, onAction, onDelete, onDuplicate, onEdit, shouldEdit]);

  return (
    <Dropdown overlay={<Menu items={menuItems} />} placement="bottomLeft" trigger={['click']}>
      <span data-cy={`test-actions-button-${id}`} className="ant-dropdown-link" onClick={e => e.stopPropagation()}>
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default ResourceCardActions;
