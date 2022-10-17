import {Modal} from 'antd';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  onConfirm(): void;
  title: string;
  heading?: string;
}

const ConfirmationModal = ({isOpen, title, heading = 'Delete Confirmation', onClose, onConfirm}: IProps) => {
  return (
    <Modal
      cancelText="Cancel"
      okText="Delete"
      onCancel={onClose}
      onOk={onConfirm}
      title={heading}
      visible={isOpen}
      data-cy="delete-confirmation-modal"
    >
      <p>{title}</p>
    </Modal>
  );
};

export default ConfirmationModal;
