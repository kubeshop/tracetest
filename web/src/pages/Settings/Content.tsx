import {Tabs} from 'antd';
import Analytics from 'components/Settings/Analytics';
import DataStore from 'components/Settings/DataStore';
import Demo from 'components/Settings/Demo';
import Linter from 'components/Settings/Linter/Linter';
import Polling from 'components/Settings/Polling';
import BetaBadge from 'components/BetaBadge/BetaBadge';
import * as S from './Settings.styled';

const TabsKeys = {
  Analytics: 'analytics',
  DataStore: 'dataStore',
  Demo: 'demo',
  Polling: 'polling',
  Linter: 'linter',
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
        <Tabs.TabPane key={TabsKeys.Demo} tab="Demo">
          <Demo />
        </Tabs.TabPane>
        <Tabs.TabPane
          key={TabsKeys.Linter}
          tab={
            <S.TabTextContainer>
              Linter
              <BetaBadge />
            </S.TabTextContainer>
          }
        >
          <Linter />
        </Tabs.TabPane>
      </Tabs>
    </S.TabsContainer>
  </S.Container>
);

export default Content;
