import {Tabs, TabsProps} from 'antd';
import {useEffect, useMemo, useState} from 'react';
import {useParams} from 'react-router-dom';
import RunDetailAutomate from 'components/RunDetailAutomate';
import RunDetailTest from 'components/RunDetailTest';
import RunDetailTrace from 'components/RunDetailTrace';
import RunDetailTrigger from 'components/RunDetailTrigger';
import {RunDetailModes} from 'constants/TestRun.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import Test from 'models/Test.model';
import {isRunStateSucceeded} from 'models/TestRun.model';
import {useNotification} from 'providers/Notification/Notification.provider';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import UserSelectors from 'selectors/User.selectors';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {ConfigMode} from 'types/DataStore.types';
import HeaderLeft from './HeaderLeft';
import HeaderRight from './HeaderRight';
import * as S from './RunDetailLayout.styled';

interface IProps {
  test: Test;
}

const renderTabBar: TabsProps['renderTabBar'] = (props, DefaultTabBar) => (
  <S.ContainerHeader data-cy="run-detail-header">
    <DefaultTabBar {...props} className="site-custom-tab-bar" />
  </S.ContainerHeader>
);

const renderTab = (title: string, testId: string, runId: number, mode: string) => (
  <S.TabLink $isActive={mode === title.toLowerCase()} to={`/test/${testId}/run/${runId}/${title.toLowerCase()}`}>
    {title}
  </S.TabLink>
);

const RunDetailLayout = ({test: {id, name, trigger}, test}: IProps) => {
  const {mode = RunDetailModes.TRIGGER} = useParams();
  const {showNotification} = useNotification();
  const {isError, run, runEvents} = useTestRun();
  const {dataStoreConfig} = useSettingsValues();
  const [prevState, setPrevState] = useState(run.state);
  useDocumentTitle(`${name} - ${run.state}`);
  const runOriginPath = useAppSelector(UserSelectors.selectRunOriginPath);

  useEffect(() => {
    const isNoTracingMode = dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

    if (isRunStateSucceeded(run.state) && !isRunStateSucceeded(prevState)) {
      showNotification({
        type: 'success',
        title: isNoTracingMode
          ? 'Response received. Skipping looking for trace as you are in No-Tracing Mode'
          : 'Trace has been fetched successfully',
      });
    }

    setPrevState(run.state);
  }, [dataStoreConfig.mode, prevState, run.state, showNotification]);

  const tabBarExtraContent = useMemo(
    () => ({
      left: <HeaderLeft name={name} triggerType={trigger.type.toUpperCase()} origin={runOriginPath} />,
      right: <HeaderRight testId={id} />,
    }),
    [id, name, trigger.type, runOriginPath]
  );

  return (
    <S.Container>
      <Tabs
        activeKey={mode}
        centered
        onChange={activeKey => {
          TestRunAnalyticsService.onChangeMode(activeKey as RunDetailModes);
        }}
        renderTabBar={renderTabBar}
        tabBarExtraContent={tabBarExtraContent}
      >
        <Tabs.TabPane tab={renderTab('Trigger', id, run.id, mode)} key={RunDetailModes.TRIGGER}>
          <RunDetailTrigger test={test} run={run} runEvents={runEvents} isError={isError} />
        </Tabs.TabPane>
        <Tabs.TabPane tab={renderTab('Trace', id, run.id, mode)} key={RunDetailModes.TRACE}>
          <RunDetailTrace run={run} runEvents={runEvents} testId={id} />
        </Tabs.TabPane>
        <Tabs.TabPane tab={renderTab('Test', id, run.id, mode)} key={RunDetailModes.TEST}>
          <RunDetailTest run={run} runEvents={runEvents} testId={id} />
        </Tabs.TabPane>
        <Tabs.TabPane tab={renderTab('Automate', id, run.id, mode)} key={RunDetailModes.AUTOMATE}>
          <RunDetailAutomate test={test} run={run} />
        </Tabs.TabPane>
      </Tabs>
    </S.Container>
  );
};

export default RunDetailLayout;
