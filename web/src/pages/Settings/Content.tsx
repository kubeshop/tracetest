import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
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
  Analyzer: 'analyzer',
};

const Content = () => {
  const [query, setQuery] = useSearchParams();

  return (
    <S.Container>
      <S.Header>
        <S.Title>Settings</S.Title>
      </S.Header>

      <S.TabsContainer>
        <Tabs
          size="small"
          defaultActiveKey={query.get('tab') || TabsKeys.DataStore}
          onChange={newTab => {
            setQuery([['tab', newTab]]);
          }}
        >
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
            key={TabsKeys.Analyzer}
            tab={
              <S.TabTextContainer>
                Analyzer
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
};

export default Content;
