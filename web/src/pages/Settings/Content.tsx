import {Tabs} from 'antd';
import Analytics from 'components/Settings/Analytics';
import DataStore from 'components/Settings/DataStore';
import Polling from 'components/Settings/Polling';
import * as S from './Settings.styled';

const TabsKeys = {
  DataStore: 'dataStore',
  Analytics: 'analytics',
  Polling: 'polling',
};

const Content = () => (
  <S.Container>
    <S.Header>
      <S.Title>Settings</S.Title>
    </S.Header>

    <S.TabsContainer>
      <Tabs size="small">
        <Tabs.TabPane key={TabsKeys.DataStore} tab="Configure Data Store">
          <DataStore />
        </Tabs.TabPane>
        <Tabs.TabPane key={TabsKeys.Analytics} tab="Analytics">
          <Analytics />
        </Tabs.TabPane>
        <Tabs.TabPane key={TabsKeys.Polling} tab="Trace Polling">
          <Polling />
        </Tabs.TabPane>
      </Tabs>
    </S.TabsContainer>
  </S.Container>
);

export default Content;
