import {Form, Tabs, TabsProps} from 'antd';
import {useEffect, useMemo, useState} from 'react';
import CreateTest from 'components/CreateTest';
import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import {TriggerTypes} from 'constants/Test.constants';
import {RunDetailModes} from 'constants/TestRun.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {TDraftTest} from 'types/Test.types';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
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
  const {onCreateTest, onUpdatePlugin, draftTest} = useCreateTest();
  useDocumentTitle(`Create - ${triggerType} test`);

  useEffect(() => {
    onUpdatePlugin(TriggerTypeToPlugin[triggerType]);
  }, [onUpdatePlugin, triggerType]);

  const plugin = TriggerTypeToPlugin[triggerType];
  const [isValid, setIsValid] = useState(false);

  const onValidate = useValidateTestDraft({pluginName: plugin.name, setIsValid});
  const [form] = Form.useForm<TDraftTest>();

  useEffect(() => {
    onValidate({}, draftTest);
  }, []);

  const tabBarExtraContent = useMemo(
    () => ({
      left: <HeaderLeft triggerType={triggerType} origin="/" />,
      right: <HeaderRight />,
    }),
    [triggerType]
  );

  return (
    <S.Container>
      <Form<TDraftTest>
        autoComplete="off"
        data-cy="edit-test-modal"
        form={form}
        layout="vertical"
        name={FORM_ID}
        initialValues={draftTest}
        onFinish={values => onCreateTest(values, plugin)}
        onValuesChange={onValidate}
      >
        <Tabs
          activeKey={RunDetailModes.TRIGGER}
          centered
          className="create-test-tabs"
          onChange={activeKey => TestRunAnalyticsService.onChangeMode(activeKey as RunDetailModes)}
          renderTabBar={renderTabBar}
          tabBarExtraContent={tabBarExtraContent}
        >
          <Tabs.TabPane tab={renderTab('Trigger', triggerType, false, true)} key={RunDetailModes.TRIGGER}>
            <CreateTest isValid={isValid} triggerType={triggerType!} />
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
