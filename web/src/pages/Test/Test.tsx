import {Tabs} from 'antd';
import Title from 'antd/lib/typography/Title';

import Trace from '../Trace';
import Assertions from './Assertions';
import * as S from './Test.styled';

const TestPage = () => {
  return (
    <>
      <S.Header>
        <Title level={3}>Test Name</Title>
      </S.Header>
      <S.Content>
        <Tabs defaultActiveKey="1">
          <Tabs.TabPane tab="Trace Details" key="1">
            <Trace />
          </Tabs.TabPane>
          <Tabs.TabPane tab="Assertions" key="2">
            <Assertions />
          </Tabs.TabPane>
        </Tabs>
      </S.Content>
    </>
  );
};

export default TestPage;
