import {Tabs} from 'antd';
import DataStore from 'components/Settings/DataStore';
import {useConfig} from 'providers/Config/Config.provider';
import * as S from './Settings.styled';

const TabsKeys = {
  DataStore: 'dataStore',
};

const Content = () => {
  const {config} = useConfig();

  return (
    <S.Container>
      <S.Header>
        <S.Title>Settings</S.Title>
      </S.Header>

      <S.TabsContainer>
        <Tabs size="small">
          <Tabs.TabPane key={TabsKeys.DataStore} tab="Configure Data Store">
            <DataStore config={config} />
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default Content;
