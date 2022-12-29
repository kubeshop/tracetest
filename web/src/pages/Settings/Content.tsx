// import {Tabs} from 'antd';
import DataStore from 'components/Settings/DataStore';
import * as S from './Settings.styled';

/* const TabsKeys = {
  DataStore: 'dataStore',
}; */

const Content = () => (
  <S.Container>
    <S.Header>
      <S.Title>Configure Data Store</S.Title>
    </S.Header>

    {/* <S.TabsContainer>
        <Tabs size="small">
          <Tabs.TabPane key={TabsKeys.DataStore} tab="Configure Data Store"> */}
    <DataStore />
    {/* </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer> */}
  </S.Container>
);

export default Content;
