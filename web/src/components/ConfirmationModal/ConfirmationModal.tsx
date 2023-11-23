import {Modal} from 'antd';
import useKeyEvent, {Keys} from 'hooks/useKeyEvent';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  onConfirm(): void;
  title: React.ReactNode;
  heading?: string;
  okText?: string;
  cancelText?: string;
}

const ConfirmationModal = ({
  isOpen,
  title,
  heading = 'Delete Confirmation',
  onClose,
  onConfirm,
  okText = 'Delete',
  cancelText = 'Cancel',
}: IProps) => {
  useKeyEvent([Keys.Enter], onConfirm);

  return (
    <Modal
      cancelText={cancelText}
      okText={okText}
      onCancel={onClose}
      onOk={onConfirm}
      title={heading}
      visible={isOpen}
      data-cy="confirmation-modal"
    >
      {typeof title === 'string' ? <p>{title}</p> : title}
    </Modal>
  );
};

export default ConfirmationModal;
