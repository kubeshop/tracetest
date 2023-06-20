import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
import CliCommand from './methods/CLICommand';
import * as S from './RunDetailAutomateMethods.styled';

const TabsKeys = {
  CLI: 'cli',
};

const RunDetailAutomateMethods = () => {
  const [query, updateQuery] = useSearchParams();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Running Techniques</S.Title> <S.Subtitle>Methods to automate the running of this test</S.Subtitle>
      </S.TitleContainer>
      <S.TabsContainer>
        <Tabs
          defaultActiveKey={query.get('tab') || TabsKeys.CLI}
          data-cy="run-detail-automate-methods"
          size="small"
          onChange={newTab => {
            updateQuery([['tab', newTab]]);
          }}
        >
          <Tabs.TabPane key={TabsKeys.CLI} tab="CLI">
            <CliCommand />
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailAutomateMethods;
