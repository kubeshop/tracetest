import {Modal} from 'antd';

interface IProps {
  cancelText?: string;
  currentVersion: number;
  description?: string;
  isOpen: boolean;
  latestVersion: number;
  okText?: string;
  onCancel(): void;
  onConfirm(): void;
}

const VersionMismatchModal = ({
  cancelText = 'Cancel',
  currentVersion,
  description,
  isOpen,
  latestVersion,
  okText = 'OK',
  onCancel,
  onConfirm,
}: IProps) => (
  <Modal
    cancelText={cancelText}
    okText={okText}
    onCancel={onCancel}
    onOk={onConfirm}
    title="Version Mismatch"
    visible={isOpen}
  >
    <p>
      You are viewing version {currentVersion} of this test, and the latest version is {latestVersion}.
    </p>
    <p>{description}</p>
  </Modal>
);

export default VersionMismatchModal;
