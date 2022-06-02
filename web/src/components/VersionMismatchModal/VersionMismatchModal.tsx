import {Modal} from 'antd';

interface IProps {
  currentVersion: number;
  description?: string;
  isOpen: boolean;
  latestVersion: number;
  onCancel(): void;
  onConfirm(): void;
}

const VersionMismatchModal = ({currentVersion, description, isOpen, latestVersion, onCancel, onConfirm}: IProps) => (
  <Modal title="Version Mismatch" visible={isOpen} onOk={onConfirm} onCancel={onCancel}>
    <p>
      You are viewing version {currentVersion} of this test, and the latest version is {latestVersion}.
    </p>
    <p>{description}</p>
  </Modal>
);

export default VersionMismatchModal;
