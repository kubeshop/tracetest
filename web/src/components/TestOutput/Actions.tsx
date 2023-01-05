import {DeleteOutlined, EditOutlined} from '@ant-design/icons';
import {Button, Tooltip} from 'antd';
import * as S from './TestOutput.styled';

interface IProps {
  isDeleted: boolean;
  onDelete(): void;
  onEdit(): void;
}

const Actions = ({isDeleted, onEdit, onDelete}: IProps) => {
  return (
    <S.ActionsContainer>
      {!isDeleted && (
        <Tooltip title="Edit output">
          <Button
            data-cy="output-item-actions-edit"
            icon={<EditOutlined />}
            onClick={event => {
              event.stopPropagation();
              onEdit();
            }}
            size="small"
            type="link"
          />
        </Tooltip>
      )}
      {!isDeleted && (
        <Tooltip title="Delete output">
          <Button
            data-cy="output-item-actions-delete"
            icon={<DeleteOutlined />}
            onClick={event => {
              event.stopPropagation();
              onDelete();
            }}
            size="small"
            type="link"
          />
        </Tooltip>
      )}
    </S.ActionsContainer>
  );
};

export default Actions;
