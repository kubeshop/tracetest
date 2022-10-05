import {Form} from 'antd';
import {useCallback} from 'react';
import PluginCard from 'components/PluginCard';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {ComponentNames, Plugins} from 'constants/Plugins.constants';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import CreateTestAnalyticsService from 'services/Analytics/CreateTestAnalytics.service';
import * as S from './SelectPlugin.styled';

const pluginList = Object.values(Plugins);

const SelectPlugin = () => {
  const {onUpdatePlugin, pluginName, onNext} = useCreateTest();
  const [form] = Form.useForm();

  const handleSubmit = useCallback(() => {
    onNext();
  }, [onNext]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>What kind of trigger do you want to use to initiate this Tracetest?</Step.Title>
        <Form id={ComponentNames.SelectPlugin} form={form} onFinish={handleSubmit}>
          <S.PluginCardList>
            {pluginList.map(plugin => (
              <PluginCard
                plugin={plugin}
                key={plugin.name}
                onSelect={selectedPlugin => {
                  CreateTestAnalyticsService.onPluginSelected(selectedPlugin.name);
                  onUpdatePlugin(selectedPlugin);
                }}
                isSelected={pluginName === plugin.name}
              />
            ))}
          </S.PluginCardList>
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default SelectPlugin;
