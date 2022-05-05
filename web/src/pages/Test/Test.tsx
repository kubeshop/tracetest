import {skipToken} from '@reduxjs/toolkit/dist/query';
import {withTracker} from 'ga-4-react';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {Button, Tabs, Typography} from 'antd';
import Title from 'antd/lib/typography/Title';
import {CloseOutlined, ArrowLeftOutlined} from '@ant-design/icons';
import {useLocation, useNavigate, useParams} from 'react-router-dom';
import {useGetTestByIdQuery, useGetResultByIdQuery, useGetResultListQuery} from 'redux/apis/Test.api';

import Trace from 'components/Trace';
import Layout from 'components/Layout';
import TestStateBadge from 'components/TestStateBadge';

import * as S from './Test.styled';
import TestDetails from './TestDetails';
import {TestState} from '../../constants/TestRunResult.constants';
import {ITestRunResult} from '../../types/TestRunResult.types';

interface TracePane {
  key: string;
  title: string;
  content: any;
}

const TestPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const {id} = useParams();
  const [tracePanes, setTracePanes] = useState<TracePane[]>([]);
  const [activeTabKey, setActiveTabKey] = useState<string>('1');
  const [readRoute, setReadRoute] = useState(false);
  const {data: test} = useGetTestByIdQuery(id as string);
  const {data: testResultList = [], isLoading} = useGetResultListQuery(id ?? skipToken, {
    pollingInterval: 5000,
  });
  const query = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const activeTestResult = testResultList.find(r => r.resultId === activeTabKey);
  const {data: activeTestResultDetails} = useGetResultByIdQuery(
    activeTestResult?.resultId && id ? {testId: id, resultId: activeTestResult.resultId} : skipToken
  );

  const isTraceLoading =
    (activeTestResult?.resultId &&
      (activeTestResultDetails?.state === TestState.AWAITING_TRACE ||
        activeTestResultDetails?.state === TestState.EXECUTING)) ||
    false;

  const handleCloseTab = useCallback(() => {
    const panes = tracePanes.filter(pane => pane.key !== activeTabKey);

    navigate(`/test/${id}`);

    setTracePanes(panes);
    setActiveTabKey('1');
  }, [activeTabKey, id, navigate, tracePanes]);

  const handleSelectTestResult = useCallback(
    (result: ITestRunResult) => {
      const itExists = Boolean(tracePanes.find(pane => pane.key === result.resultId));

      if (!itExists) {
        const newTabIndex = testResultList.findIndex(r => r.resultId === result.resultId) + 1;
        const tracePane = {
          key: result.resultId,
          title: `Trace #${newTabIndex}`,
          content: (
            <Trace
              testId={id!}
              testResultId={result.resultId}
              onDismissTrace={handleCloseTab}
              onRunTest={handleSelectTestResult}
            />
          ),
        };

        setTracePanes([...tracePanes, tracePane]);
      }

      navigate(`/test/${id}?resultId=${result.resultId}`);
      setActiveTabKey(`${result.resultId}`);
    },
    [handleCloseTab, id, navigate, testResultList, tracePanes]
  );

  useEffect(() => {
    const resultId = query.get('resultId');

    if (test && resultId && resultId !== activeTabKey && !readRoute && testResultList.length) {
      const testResult = testResultList.find(({resultId: rId}) => rId === resultId);

      if (testResult) handleSelectTestResult(testResult);
      setReadRoute(true);
    }
  }, [location, test, testResultList, readRoute]);

  const onChangeTab = useCallback(
    (tabKey: string) => {
      if (Number.isNaN(Number(tabKey))) navigate(`/test/${id}?resultId=${tabKey}`);
      else navigate(`/test/${id}`);

      setActiveTabKey(tabKey);
    },
    [id, navigate]
  );

  const onEditTab = useCallback(
    (targetKey: any) => {
      let activeKey = activeTabKey;
      let lastIndex = -1;
      tracePanes.forEach((pane, i) => {
        if (pane.key === targetKey) {
          lastIndex = i - 1;
        }
      });
      const panes = tracePanes.filter(pane => pane.key !== targetKey);
      if (tracePanes.length && activeTabKey === targetKey) {
        if (lastIndex >= 0) {
          activeKey = panes[lastIndex].key;
        } else {
          activeKey = '1';
        }
      }

      setTracePanes(panes);
      setActiveTabKey(activeKey);
    },
    [activeTabKey, tracePanes]
  );

  return (
    <Layout>
      <ReactFlowProvider>
        <S.TestTabs
          loading={isTraceLoading.toString()}
          tabBarExtraContent={{
            left: (
              <S.Header>
                <Button type="text" shape="circle" onClick={() => navigate('/')}>
                  <ArrowLeftOutlined style={{fontSize: 24, marginRight: 16}} />
                </Button>
                <Title style={{margin: 0}} level={3}>
                  {test?.name}
                </Title>
              </S.Header>
            ),
            right: activeTestResult?.resultId && activeTestResultDetails?.state && (
              <div style={{marginRight: 24}}>
                <Typography.Text style={{marginRight: 8, color: '#8C8C8C', fontSize: 14}}>Test status:</Typography.Text>
                <TestStateBadge style={{fontSize: 16}} testState={activeTestResultDetails?.state} />
              </div>
            ),
          }}
          hideAdd
          defaultActiveKey="1"
          activeKey={activeTabKey}
          onChange={onChangeTab}
          type="editable-card"
          onEdit={onEditTab}
          style={{flexGrow: 1, display: 'flex', margin: 0}}
        >
          <Tabs.TabPane tab="Test Details" key="1" closeIcon={<CloseOutlined hidden />}>
            <S.Wrapper>
              <TestDetails
                testResultList={testResultList}
                isLoading={isLoading}
                testId={id!}
                url={test?.serviceUnderTest.request.url}
                onSelectResult={handleSelectTestResult}
              />
            </S.Wrapper>
          </Tabs.TabPane>
          {tracePanes.map(item => (
            <Tabs.TabPane tab={item.title} key={item.key}>
              <S.Wrapper>{item.content}</S.Wrapper>
            </Tabs.TabPane>
          ))}
        </S.TestTabs>
      </ReactFlowProvider>
    </Layout>
  );
};

export default withTracker(TestPage);
