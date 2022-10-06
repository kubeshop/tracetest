import {Tabs, TabsProps} from 'antd';

import RunDetailTest from 'components/RunDetailTest';
import RunDetailTrace from 'components/RunDetailTrace';
import RunDetailTrigger from 'components/RunDetailTrigger';
import {RunDetailModes} from 'constants/TestRun.constants';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useMemo} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {TTest} from 'types/Test.types';
import {Steps} from '../GuidedTour/traceStepList';
import HeaderLeft from './HeaderLeft';
import HeaderRight from './HeaderRight';
import * as S from './RunDetailLayout.styled';

interface IProps {
  test: TTest;
}

const renderTabBar: TabsProps['renderTabBar'] = (props, DefaultTabBar) => (
  <S.ContainerHeader
    data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Switcher)}
    data-cy="run-detail-header"
  >
    <DefaultTabBar {...props} className="site-custom-tab-bar" />
  </S.ContainerHeader>
);

const RunDetailLayout = ({test: {id, name, trigger, version = 1}, test}: IProps) => {
  const navigate = useNavigate();
  const {mode = RunDetailModes.TRIGGER} = useParams();
  const {isError, run} = useTestRun();

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
          navigate(`/test/${id}/run/${run.id}/${activeKey}`);
        }}
        renderTabBar={renderTabBar}
        tabBarExtraContent={tabBarExtraContent}
        destroyInactiveTabPane
      >
        <Tabs.TabPane tab="Trigger" key={RunDetailModes.TRIGGER}>
          <RunDetailTrigger test={test} run={run} isError={isError} />
        </Tabs.TabPane>
        <Tabs.TabPane tab="Trace" key={RunDetailModes.TRACE}>
          <RunDetailTrace run={run} testId={id} />
        </Tabs.TabPane>
        <Tabs.TabPane tab="Test" key={RunDetailModes.TEST}>
          <RunDetailTest run={run} testId={id} />
        </Tabs.TabPane>
      </Tabs>
    </S.Container>
  );
};

export default RunDetailLayout;
