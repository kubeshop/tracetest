import {Tabs} from 'antd';
import Title from 'antd/lib/typography/Title';
import {useParams} from 'react-router-dom';
import {useGetTestByIdQuery} from '../../services/TestService';

import Trace from '../Trace';
import Assertions from './Assertions';
import * as S from './Test.styled';

const TestPage = () => {
  const {id} = useParams();
  const {data: test} = useGetTestByIdQuery(id as string);
  return (
    <>
      <S.Header>
        <Title level={3}>{test?.name}</Title>
      </S.Header>
      <S.Content>
        <Tabs defaultActiveKey="1">
          {test && (
            <Tabs.TabPane tab="Trace Details" key="1">
              <Trace test={test} />
            </Tabs.TabPane>
          )}
          <Tabs.TabPane tab="Assertions" key="2">
            <Assertions />
          </Tabs.TabPane>
        </Tabs>
      </S.Content>
    </>
  );
};

export default TestPage;
