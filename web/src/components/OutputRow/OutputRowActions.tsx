import {Dropdown, Menu} from 'antd';
import {useMemo} from 'react';
import * as S from './OutputRow.styled';

interface IProps {
  name: string;
  onDelete(): void;
  onEdit(): void;
}

const OutputRowActions = ({name, onEdit, onDelete}: IProps) => {
  const items = useMemo(
    () => [
      {key: 'edit', label: <span data-cy="output-item-actions-edit">Edit</span>, onClick: onEdit},
      {key: 'delete', label: <span data-cy="output-item-actions-delete">Delete</span>, onClick: onDelete},
    ],
    [onDelete, onEdit]
  );

  return (
    <Dropdown overlay={<Menu items={items} />} placement="bottomLeft" trigger={['click']}>
      <span data-cy={`output-actions-button-${name}`} className="ant-dropdown-link" onClick={e => e.stopPropagation()}>
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default OutputRowActions;
