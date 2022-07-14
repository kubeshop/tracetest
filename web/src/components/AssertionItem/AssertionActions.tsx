import {DeleteOutlined, EditOutlined, SelectOutlined, UndoOutlined} from '@ant-design/icons';
import {Button, Tooltip} from 'antd';

import * as S from './AssertionItem.styled';

interface IProps {
  isDeleted: boolean;
  isDraft: boolean;
  isSelected: boolean;
  onDelete(): void;
  onEdit(): void;
  onRevert(): void;
  onSelect(): void;
}

const AssertionActions = ({isDeleted, isDraft, isSelected, onDelete, onEdit, onRevert, onSelect}: IProps) => (
  <>
    {isDraft && <S.ActionTag>draft</S.ActionTag>}
    {isDeleted && <S.ActionTag>deleted</S.ActionTag>}
    {isDraft && (
      <Tooltip title="Revert assertion">
        <Button
          icon={<UndoOutlined />}
          onClick={event => {
            event.stopPropagation();
            onRevert();
          }}
          size="small"
          type="link"
        />
      </Tooltip>
    )}
    <Tooltip title="Edit assertion">
      <Button
        data-cy="edit-assertion-button"
        icon={<EditOutlined />}
        onClick={event => {
          event.stopPropagation();
          onEdit();
        }}
        size="small"
        type="link"
      />
    </Tooltip>
    <Tooltip title="Delete assertion">
      <Button
        icon={<DeleteOutlined />}
        onClick={event => {
          event.stopPropagation();
          onDelete();
        }}
        size="small"
        type="link"
      />
    </Tooltip>
    <Tooltip title={isSelected ? 'Unselect assertion' : 'Select assertion'}>
      <Button
        icon={<SelectOutlined />}
        onClick={event => {
          event.stopPropagation();
          onSelect();
        }}
        size="small"
        type="link"
      />
    </Tooltip>
  </>
);

export default AssertionActions;
