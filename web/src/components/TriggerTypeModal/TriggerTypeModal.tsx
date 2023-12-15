import {CreateTriggerTypeToPlugin} from 'constants/Plugins.constants';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import CreateTestAnalyticsService from 'services/Analytics/CreateTestAnalytics.service';
import TriggerTypeCard from './TriggerTypeCard';
import * as S from './TriggerTypeModal.styled';
import { withCustomization } from '../../providers/Customization';

const pluginList = Object.values(CreateTriggerTypeToPlugin);

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

      <S.CardList>
        {pluginList.map(plugin => (
          <TriggerTypeCard
            key={plugin.name}
            onClick={selectedPlugin => {
              CreateTestAnalyticsService.onPluginSelected(selectedPlugin.name);
              navigate(`test/create/${selectedPlugin.type}`);
            }}
            plugin={plugin}
          />
        ))}
      </S.CardList>
    </S.Modal>
  );
};

export default withCustomization(TriggerTypeModal, 'triggerTypeModal');
