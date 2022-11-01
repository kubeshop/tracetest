import {Button} from 'antd';

import * as S from './OutputModal.styled';

interface IProps {
  isEditing: boolean;
  isValid: boolean;
  onCancel(): void;
  onSave(): void;
}

const OutputModalFooter = ({isEditing, isValid, onCancel, onSave}: IProps) => (
  <S.Footer>
    <span>
      <Button data-cy="output-modal-cancel-button" type="primary" ghost onClick={onCancel}>
        Cancel
      </Button>
    </span>
    <span>
      <Button htmlType="submit" data-cy="output-modal-save-button" disabled={!isValid} type="primary" onClick={onSave}>
        {isEditing ? 'Update' : 'Add'}
      </Button>
    </span>
  </S.Footer>
);

export default OutputModalFooter;
