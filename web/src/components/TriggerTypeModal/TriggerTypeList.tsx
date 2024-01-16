import {useState} from 'react';
import CreateTestAnalytics from 'services/Analytics/CreateTest.service';
import {CreateTriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {withCustomization} from 'providers/Customization';
import TriggerTypeCard from './TriggerTypeCard';
import * as S from './TriggerTypeModal.styled';

const pluginList = Object.values(CreateTriggerTypeToPlugin);

interface IProps {
  onClick(plugin: TriggerTypes): void;
}

const TriggerTypeList = ({onClick}: IProps) => {
  const [selectedType, setSelectedType] = useState<TriggerTypes | undefined>(undefined);

  return (
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
  );
};

export default withCustomization(TriggerTypeList, 'TriggerTypeList');
