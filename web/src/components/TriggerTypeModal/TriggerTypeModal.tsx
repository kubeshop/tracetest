import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import TriggerTypeList from './TriggerTypeList';
import * as S from './TriggerTypeModal.styled';

interface IProps {
  isOpen: boolean;
  onClose(): void;
}

const TriggerTypeModal = ({isOpen, onClose}: IProps) => {
  const {navigate} = useDashboard();

  return (
    <S.Modal
      onCancel={onClose}
      footer={null}
      title={<S.Title level={2}>Create a new test</S.Title>}
      visible={isOpen}
      width={625}
    >
      <S.Title level={3} $marginBottom={16}>
        What kind of trigger do you want to use to initiate this Tracetest?
      </S.Title>

      <TriggerTypeList onClick={type => navigate(`test/create/${type}`)} />
    </S.Modal>
  );
};

export default TriggerTypeModal;
