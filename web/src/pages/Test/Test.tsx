import {skipToken} from '@reduxjs/toolkit/dist/query';
import {useCallback, useEffect, useMemo, useState} from 'react';
import {ReactFlowProvider} from 'react-flow-renderer';
import {Button, Tabs} from 'antd';
import Title from 'antd/lib/typography/Title';
import {CloseOutlined, ArrowLeftOutlined} from '@ant-design/icons';
import {useLocation, useNavigate, useParams} from 'react-router-dom';

import {ITestResult} from 'types';
import {useGetTestByIdQuery, useGetTestResultsQuery} from 'services/TestService';
import Trace from 'components/Trace';
import Layout from 'components/Layout';

import Assertions from './Assertions';
import * as S from './Test.styled';
import TestDetails from './TestDetails';

interface ITestRouteState {
  testRun: ITestResult;
}
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
  const {data: test} = useGetTestByIdQuery(id as string);
  const {data: testResultList = [], isLoading} = useGetTestResultsQuery(id ?? skipToken);
  const query = useMemo(() => new URLSearchParams(location.search), [location.search]);

  const handleSelectTestResult = useCallback(
    (result: ITestResult) => {
      const itExists = Boolean(tracePanes.find(pane => pane.key === result.resultId));

      if (!itExists) {
        const newTabIndex = testResultList.findIndex(r => r.resultId === result.resultId) + 1;
        const tracePane = {
          key: result.resultId,
          title: `Trace #${newTabIndex}`,
          content: <Trace test={test!} testResultId={result.resultId} />,
        };

        setTracePanes([...tracePanes, tracePane]);
      }

      navigate(`/test/${id}?resultId=${result.resultId}`);
      setActiveTabKey(`${result.resultId}`);
    },
    [id, navigate, test, testResultList, tracePanes]
  );

  useEffect(() => {
    if ((location?.state as ITestRouteState)?.testRun && test) {
      handleSelectTestResult((location.state as ITestRouteState).testRun);
    }
  }, [location, test]);

  useEffect(() => {
    const resultId = query.get('resultId');

    if (test && resultId && resultId !== activeTabKey) {
      const testResult = testResultList.find(({resultId: rId}) => rId === resultId);

      if (testResult) handleSelectTestResult(testResult);
    }
  }, [location, test, testResultList]);

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
        <Tabs
          tabBarExtraContent={{
            left: (
              <S.Header>
                <Button type="text" shape="circle" onClick={() => navigate(-1)}>
                  <ArrowLeftOutlined style={{fontSize: 24, marginRight: 16}} />
                </Button>
                <Title style={{margin: 0}} level={3}>
                  {test?.name}
                </Title>
              </S.Header>
            ),
          }}
          hideAdd
          defaultActiveKey="1"
          activeKey={activeTabKey}
          onChange={onChangeTab}
          type="editable-card"
          onEdit={onEditTab}
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
          {Boolean(test?.assertions?.length) && (
            <Tabs.TabPane tab="Test Assertions" key="2" closeIcon={<CloseOutlined hidden />}>
              <S.Wrapper>
                <Assertions />
              </S.Wrapper>
            </Tabs.TabPane>
          )}
          {tracePanes.map(item => (
            <Tabs.TabPane tab={item.title} key={item.key}>
              <S.Wrapper>{item.content}</S.Wrapper>
            </Tabs.TabPane>
          ))}
        </Tabs>
      </ReactFlowProvider>
    </Layout>
  );
};

export default TestPage;
