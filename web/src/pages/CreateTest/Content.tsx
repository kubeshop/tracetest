import {Tabs, TabsProps} from 'antd';
// import {useMemo} from 'react';
// import HeaderLeft from 'components/RunDetailLayout/HeaderLeft';
// import HeaderRight from 'components/RunDetailLayout/HeaderRight';
import CreateTest from 'components/CreateTest';
import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import {TriggerTypes} from 'constants/Test.constants';
import {RunDetailModes} from 'constants/TestRun.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';

interface IProps {
  triggerType: TriggerTypes;
}

const renderTabBar: TabsProps['renderTabBar'] = (props, DefaultTabBar) => (
  <S.ContainerHeader data-cy="run-detail-header">
    <DefaultTabBar {...props} className="site-custom-tab-bar" />
  </S.ContainerHeader>
);

const renderTab = (title: string, triggerType: TriggerTypes, isDisabled: boolean, isActive: boolean = false) => (
  <S.TabLink to={isDisabled ? '#' : `/test/create/${triggerType}`} $isActive={isActive} $isDisabled={isDisabled}>
    {title}
  </S.TabLink>
);

const Content = ({triggerType}: IProps) => {
  useDocumentTitle(`Create - ${triggerType} test`);

  /* const tabBarExtraContent = useMemo(
    () => ({
      left: <HeaderLeft name={name} triggerType={trigger.type.toUpperCase()} origin={runOriginPath} />,
      right: <HeaderRight testId={id} />,
    }),
    [id, name, trigger.type, runOriginPath]
  ); */

  return (
    <S.Container>
      <Tabs
        activeKey={RunDetailModes.TRIGGER}
        centered
        className="run-tabs"
        onChange={activeKey => {
          TestRunAnalyticsService.onChangeMode(activeKey as RunDetailModes);
        }}
        renderTabBar={renderTabBar}
        // tabBarExtraContent={tabBarExtraContent}
      >
        <Tabs.TabPane tab={renderTab('Trigger', triggerType, false, true)} key={RunDetailModes.TRIGGER}>
          <CreateTest triggerType={triggerType!} />
        </Tabs.TabPane>
        <Tabs.TabPane disabled tab={renderTab('Trace', triggerType, true)} key={RunDetailModes.TRACE} />
        <Tabs.TabPane disabled tab={renderTab('Test', triggerType, true)} key={RunDetailModes.TEST} />
        <Tabs.TabPane disabled tab={renderTab('Automate', triggerType, true)} key={RunDetailModes.AUTOMATE} />
      </Tabs>
    </S.Container>
  );
};

export default Content;
