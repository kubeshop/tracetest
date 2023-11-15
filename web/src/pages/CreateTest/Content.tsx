import {Form, Tabs, TabsProps} from 'antd';
import CreateTest from 'components/CreateTest';
import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import {getDemoByPluginMap} from 'constants/Demo.constants';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {RunDetailModes} from 'constants/TestRun.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {useCreateTest} from 'providers/CreateTest';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {useEffect, useMemo, useState} from 'react';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {TDraftTest} from 'types/Test.types';
import HeaderLeft from './HeaderLeft';
import HeaderRight from './HeaderRight';

interface IProps {
  triggerType: TriggerTypes;
}

const renderTabBar: TabsProps['renderTabBar'] = (props, DefaultTabBar) => (
  <S.ContainerHeader data-cy="create-test-header">
    <DefaultTabBar {...props} className="site-custom-tab-bar" />
  </S.ContainerHeader>
);

const renderTab = (title: string, triggerType: TriggerTypes, isDisabled: boolean, isActive: boolean = false) => (
  <S.TabLink to={isDisabled ? '#' : `/test/create/${triggerType}`} $isActive={isActive} $isDisabled={isDisabled}>
    {title}
  </S.TabLink>
);

export const FORM_ID = 'create-test';

const Content = ({triggerType}: IProps) => {
  const {onCreateTest, isLoading, initialValues} = useCreateTest();
  useDocumentTitle(`Create ${triggerType} test`);

  const [form] = Form.useForm<TDraftTest>();

  const plugin = TriggerTypeToPlugin[triggerType];
  const [isValid, setIsValid] = useState(false);
  const onValidateTest = useValidateTestDraft({pluginName: plugin.name, setIsValid});

  const {demos} = useSettingsValues();
  const demosByTriggerType = useMemo(() => {
    const demoByPluginMap = getDemoByPluginMap(demos);
    return demoByPluginMap[plugin.name];
  }, [demos, plugin.name]);

  useEffect(() => {
    onValidateTest({}, initialValues);
  }, []);

  const tabBarExtraContent = useMemo(
    () => ({
      left: <HeaderLeft triggerType={triggerType} origin="/" />,
      right: <HeaderRight demos={demosByTriggerType} />,
    }),
    [demosByTriggerType, triggerType]
  );

  return (
    <S.Container>
      <Form<TDraftTest>
        autoComplete="off"
        data-cy="create-test"
        form={form}
        layout="vertical"
        name={FORM_ID}
        initialValues={initialValues}
        onFinish={values => onCreateTest(values, plugin)}
        onValuesChange={onValidateTest}
      >
        <Tabs
          activeKey={RunDetailModes.TRIGGER}
          centered
          onChange={activeKey => TestRunAnalyticsService.onChangeMode(activeKey as RunDetailModes)}
          renderTabBar={renderTabBar}
          tabBarExtraContent={tabBarExtraContent}
        >
          <Tabs.TabPane tab={renderTab('Trigger', triggerType, false, true)} key={RunDetailModes.TRIGGER}>
            <CreateTest isLoading={isLoading} isValid={isValid} triggerType={triggerType!} />
          </Tabs.TabPane>
          <Tabs.TabPane disabled tab={renderTab('Trace', triggerType, true)} key={RunDetailModes.TRACE} />
          <Tabs.TabPane disabled tab={renderTab('Test', triggerType, true)} key={RunDetailModes.TEST} />
          <Tabs.TabPane disabled tab={renderTab('Automate', triggerType, true)} key={RunDetailModes.AUTOMATE} />
        </Tabs>
      </Form>
    </S.Container>
  );
};

export default Content;
