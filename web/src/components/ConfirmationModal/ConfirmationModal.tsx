import {Modal} from 'antd';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  onConfirm(): void;
  title: string;
}

const ConfirmationModal = ({isOpen, title, onClose, onConfirm}: IProps) => {
  return (
    <Modal
      cancelText="Cancel"
      okText="Delete"
      onCancel={onClose}
      onOk={onConfirm}
      title="Delete Confirmation"
      visible={isOpen}
      data-cy="delete-confirmation-modal"
    >
      <p>{title}</p>
    </Modal>
  );
};

export default ConfirmationModal;
