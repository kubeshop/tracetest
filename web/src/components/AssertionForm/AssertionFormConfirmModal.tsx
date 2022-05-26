import {Modal} from 'antd';

interface IProps {
  isOpen: boolean;
  onCancel(): void;
  onConfirm(): void;
  latestVersion: number;
  currentVersion: number;
}

const AssertionFormConfirmModal: React.FC<IProps> = ({isOpen, onCancel, onConfirm, latestVersion, currentVersion}) => {
  return (
    <Modal title="Version Mismatch" visible={isOpen} onOk={onConfirm} onCancel={onCancel}>
      <p>
        You are viewing version {currentVersion} of this test, and the latest version is {latestVersion}.
      </p>
      <p>Changing and saving changes will result in a new version that will become the latest</p>
    </Modal>
  );
};

export default AssertionFormConfirmModal;
