import {useState} from 'react';
import CreateTestAnalytics from 'services/Analytics/CreateTest.service';
import {Integrations} from 'constants/Integrations.constants';
import {CreateTriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {withCustomization} from 'providers/Customization';
import IntegrationTypeCard from './IntegrationTypeCard';
import TriggerTypeCard from './TriggerTypeCard';
import * as S from './TriggerTypeModal.styled';

const pluginList = Object.values(CreateTriggerTypeToPlugin);
const integrationList = Object.values(Integrations);

interface IProps {
  onClick(plugin: TriggerTypes): void;
}

const TriggerTypeList = ({onClick}: IProps) => {
  const [selectedType, setSelectedType] = useState<TriggerTypes | undefined>(undefined);

  return (
    <>
      <S.CardList>
        {pluginList.map(plugin => (
          <TriggerTypeCard
            key={plugin.name}
            isSelected={selectedType === plugin.type}
            onClick={selectedPlugin => {
              CreateTestAnalytics.onTriggerSelect(selectedPlugin.type);
              onClick(selectedPlugin.type);
              setSelectedType(selectedPlugin.type);
            }}
            plugin={plugin}
          />
        ))}
      </S.CardList>

      <S.Divider />

      <S.Title level={3} $marginBottom={16}>
        Or trigger a Tracetest via an external integration
      </S.Title>

      <S.IntegrationCardList>
        {integrationList.map(integration => (
          <IntegrationTypeCard key={integration.name} integration={integration} />
        ))}
      </S.IntegrationCardList>

      <S.Text type="secondary">* this capability is limited to the commercial version of Tracetest.</S.Text>
    </>
  );
};

export default withCustomization(TriggerTypeList, 'TriggerTypeList');
