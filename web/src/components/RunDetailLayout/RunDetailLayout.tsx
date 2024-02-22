import {Tabs, TabsProps} from 'antd';
import {useMemo} from 'react';
import {useParams} from 'react-router-dom';
import RunDetailAutomate from 'components/RunDetailAutomate';
import RunDetailTest from 'components/RunDetailTest';
import RunDetailTrace from 'components/RunDetailTrace';
import RunDetailTrigger from 'components/RunDetailTrigger';
import {RunDetailModes} from 'constants/TestRun.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import Test from 'models/Test.model';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import UserSelectors from 'selectors/User.selectors';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import HeaderLeft from './HeaderLeft';
import HeaderRight from './HeaderRight';
import * as S from './RunDetailLayout.styled';
import useRunCompletion from './hooks/useRunCompletion';

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

const RunDetailLayout = ({test: {id, name, trigger, skipTraceCollection}, test}: IProps) => {
  const {mode = RunDetailModes.TEST} = useParams();
  const {isError, run, runEvents} = useTestRun();
  useDocumentTitle(`${name} - ${run.state}`);
  const runOriginPath = useAppSelector(UserSelectors.selectRunOriginPath);

  useRunCompletion();

  const tabBarExtraContent = useMemo(
    () => ({
      left: <HeaderLeft name={name} triggerType={trigger.type.toUpperCase()} origin={runOriginPath} />,
      right: <HeaderRight testId={id} triggerType={trigger.type} />,
    }),
    [id, name, trigger.type, runOriginPath]
  );

  return (
    <S.Container>
      <Tabs
        activeKey={mode}
        centered
        className="run-tabs"
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
          <RunDetailTrace run={run} runEvents={runEvents} testId={id} skipTraceCollection={skipTraceCollection} />
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
