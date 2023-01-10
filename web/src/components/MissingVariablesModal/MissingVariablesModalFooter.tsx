import {Button} from 'antd';
import * as S from './MissingVariablesModal.styled';

interface IProps {
  isValid: boolean;
  onCancel(): void;
  onSave(): void;
}

const MissingVariablesModalFooter = ({isValid, onCancel, onSave}: IProps) => (
  <S.Footer>
    <span>
      <Button data-cy="output-modal-cancel-button" type="primary" ghost onClick={onCancel}>
        Cancel
      </Button>
    </span>
    <span>
      <Button htmlType="submit" data-cy="output-save-button" disabled={!isValid} type="primary" onClick={onSave}>
        Run
      </Button>
    </span>
  </S.Footer>
);

export default MissingVariablesModalFooter;
