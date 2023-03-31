import {Tabs, TabsProps} from 'antd';
import {useEffect, useMemo, useState} from 'react';
import {useParams} from 'react-router-dom';
import RunDetailTest from 'components/RunDetailTest';
import RunDetailTrace from 'components/RunDetailTrace';
import RunDetailTrigger from 'components/RunDetailTrigger';
import {RunDetailModes, TestState} from 'constants/TestRun.constants';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import useDocumentTitle from 'hooks/useDocumentTitle';
import Test from 'models/Test.model';
import {useNotification} from 'providers/Notification/Notification.provider';
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

const renderTab = (title: string, testId: string, runId: string, mode: string) => (
  <S.TabLink $isActive={mode === title.toLowerCase()} to={`/test/${testId}/run/${runId}/${title.toLowerCase()}`}>
    {title}
  </S.TabLink>
);

const RunDetailLayout = ({test: {id, name, trigger, version = 1}, test}: IProps) => {
  const {mode = RunDetailModes.TRIGGER} = useParams();
  const {showNotification} = useNotification();
  const {isError, run, runEvents} = useTestRun();
  const [prevState, setPrevState] = useState(run.state);
  useDocumentTitle(`${name} - ${run.state}`);

  useEffect(() => {
    if (run.state === TestState.FINISHED && prevState !== TestState.FINISHED) {
      showNotification({
        type: 'success',
        title: 'Trace has been fetched successfully',
      });
    }

    setPrevState(run.state);
  }, [prevState, run.state, showNotification]);

  const tabBarExtraContent = useMemo(
    () => ({
      left: <HeaderLeft testId={id} name={name} triggerType={trigger.type.toUpperCase()} />,
      right: <HeaderRight testId={id} testVersion={version} />,
    }),
    [id, name, trigger.type, version]
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
        destroyInactiveTabPane
      >
        <Tabs.TabPane tab={renderTab('Trigger', id, run.id, mode)} key={RunDetailModes.TRIGGER}>
          <RunDetailTrigger test={test} run={run} runEvents={runEvents} isError={isError} />
        </Tabs.TabPane>
        <Tabs.TabPane tab={renderTab('Trace', id, run.id, mode)} key={RunDetailModes.TRACE}>
          <RunDetailTrace run={run} testId={id} />
        </Tabs.TabPane>
        <Tabs.TabPane tab={renderTab('Test', id, run.id, mode)} key={RunDetailModes.TEST}>
          <RunDetailTest run={run} testId={id} />
        </Tabs.TabPane>
      </Tabs>
    </S.Container>
  );
};

export default RunDetailLayout;
