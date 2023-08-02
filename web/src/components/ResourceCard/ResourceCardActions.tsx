import {Dropdown, Menu} from 'antd';
import {useCallback, useMemo} from 'react';

import * as S from './ResourceCard.styled';

interface IProps {
  id: string;
  shouldEdit: boolean;
  onDelete(): void;
  onEdit(): void;
}

const ResourceCardActions = ({id, shouldEdit = true, onDelete, onEdit}: IProps) => {
  const onDeleteClick = useCallback(
    ({domEvent}: {domEvent: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>}) => {
      domEvent?.stopPropagation();
      onDelete();
    },
    [onDelete]
  );

  const onEditClick = useCallback(
    ({domEvent}: {domEvent: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>}) => {
      domEvent?.stopPropagation();
      onEdit();
    },
    [onEdit]
  );

  const menuItems = useMemo(() => {
    const defaultItems = [
      {key: 'delete', label: <span data-cy="test-card-delete">Delete</span>, onClick: onDeleteClick},
    ];

    return shouldEdit
      ? [{key: 'edit', label: <span data-cy="test-card-edit">Edit</span>, onClick: onEditClick}, ...defaultItems]
      : defaultItems;
  }, [onDeleteClick, onEditClick, shouldEdit]);

  return (
    <Dropdown overlay={<Menu items={menuItems} />} placement="bottomLeft" trigger={['click']}>
      <span data-cy={`test-actions-button-${id}`} className="ant-dropdown-link" onClick={e => e.stopPropagation()}>
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default ResourceCardActions;
