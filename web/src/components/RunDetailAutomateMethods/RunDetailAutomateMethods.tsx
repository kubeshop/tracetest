import {Tabs} from 'antd';
import {useSearchParams} from 'react-router-dom';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import CliCommand from './methods/CLICommand';
import * as S from './RunDetailAutomateMethods.styled';
import DeepLink from './methods/DeepLink/DeepLink';

const TabsKeys = {
  CLI: 'cli',
  DeepLink: 'deeplink',
};

export interface IMethodProps {
  environmentId?: string;
  test: Test;
  run: TestRun;
  fileName?: string;
}

const RunDetailAutomateMethods = ({test, run, fileName}: IMethodProps) => {
  const [query, updateQuery] = useSearchParams();
  const {selectedEnvironment: {id: environmentId} = {}} = useEnvironment();

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
            <CliCommand test={test} environmentId={environmentId} run={run} fileName={fileName} />
          </Tabs.TabPane>
          <Tabs.TabPane key={TabsKeys.DeepLink} tab="Deep Link">
            <DeepLink test={test} environmentId={environmentId} run={run} />
          </Tabs.TabPane>
        </Tabs>
      </S.TabsContainer>
    </S.Container>
  );
};

export default RunDetailAutomateMethods;
