import {useRef, useState} from 'react';
import {Button, Tabs} from 'antd';
import Title from 'antd/lib/typography/Title';
import {CloseOutlined} from '@ant-design/icons';
import {useParams} from 'react-router-dom';

import {ITestResult} from 'types';
import {useGetTestByIdQuery, useRunTestMutation} from 'services/TestService';
import Trace from 'components/Trace';

import Assertions from './Assertions';
import * as S from './Test.styled';
import TestDetails from './TestDetails';

interface TracePane {
  key: string;
  title: string;
  content: any;
}

const TestPage = () => {
  const {id} = useParams();
  const newTabIndexRef = useRef<number>(2);
  const [tracePanes, setTracePanes] = useState<TracePane[]>([]);
  const [activeTabKey, setActiveTabKey] = useState<string>('1');
  const {data: test} = useGetTestByIdQuery(id as string);
  const [runTest] = useRunTestMutation();

  const handleRunTest = () => {
    if (test?.id) {
      runTest(test.id);
    }
  };

  const handleSelectTestResult = (result: ITestResult) => {
    newTabIndexRef.current += 1;
    const newTabIndex = newTabIndexRef.current;
    const tracePane = {
      key: `${newTabIndex}`,
      title: `Trace ${newTabIndex}`,
      content: <Trace test={test!} testResultId={result.id} />,
    };
    setTracePanes([...tracePanes, tracePane]);
    setActiveTabKey(`${newTabIndex}`);
  };

  const onChangeTab = (tabKey: string) => {
    setActiveTabKey(tabKey);
  };

  const onEditTab = (targetKey: any) => {
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
  };

  return (
    <>
      <S.Header>
        <Title style={{margin: 0}} level={3}>
          {test?.name}
        </Title>
        <Button onClick={handleRunTest}>Generate Trace</Button>
      </S.Header>
      <S.Content>
        <Tabs
          hideAdd
          defaultActiveKey="1"
          activeKey={activeTabKey}
          onChange={onChangeTab}
          type="editable-card"
          onEdit={onEditTab}
        >
          <Tabs.TabPane tab="Test Details" key="1" closeIcon={<CloseOutlined hidden />}>
            <TestDetails testId={id!} onSelectResult={handleSelectTestResult} />
          </Tabs.TabPane>
          <Tabs.TabPane tab="Assertions" key="2" closeIcon={<CloseOutlined hidden />}>
            <Assertions />
          </Tabs.TabPane>
          {tracePanes.map(item => (
            <Tabs.TabPane tab={item.title} key={item.key}>
              {item.content}
            </Tabs.TabPane>
          ))}
        </Tabs>
      </S.Content>
    </>
  );
};

export default TestPage;
