import {Button} from 'antd';
import * as S from './OutputModal.styled';

interface IProps {
  onCancel(): void;
  onSave(): void;
  isLoading: boolean;
  isValid: boolean;
  isEditing: boolean;
}

const OutputModalFooter = ({isValid, onCancel, isLoading, isEditing, onSave}: IProps) => {
  return (
    <S.Footer>
      <span>
        <Button data-cy="output-modal-cancel-button" type="primary" ghost onClick={onCancel}>
          Cancel
        </Button>
      </span>
      <span>
        <Button
          htmlType="submit"
          data-cy="output-modal-save-button"
          disabled={!isValid}
          type="primary"
          loading={isLoading}
          onClick={onSave}
        >
          {isEditing ? 'Save' : 'Create'}
        </Button>
      </span>
    </S.Footer>
  );
};

export default OutputModalFooter;
