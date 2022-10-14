import {useMemo} from 'react';
import {Dropdown, Menu} from 'antd';
import * as S from './OutputRow.styled';

interface IProps {
  outputId: string;
  onDelete(): void;
  onEdit(): void;
}

const OutputRowActions: React.FC<IProps> = ({outputId, onEdit, onDelete}) => {
  const items = useMemo(
    () => [
      {key: 'edit', label: <span data-cy="output-row-actions-edit">Edit</span>, onClick: onEdit},
      {key: 'delete', label: <span data-cy="output-row-actions-delete">Delete</span>, onClick: onDelete},
    ],
    [onDelete, onEdit]
  );

  return (
    <Dropdown overlay={<Menu items={items} />} placement="bottomLeft" trigger={['click']}>
      <span
        data-cy={`test-actions-button-${outputId}`}
        className="ant-dropdown-link"
        onClick={e => e.stopPropagation()}
      >
        <S.ActionButton />
      </span>
    </Dropdown>
  );
};

export default OutputRowActions;
