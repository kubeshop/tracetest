import {snakeCase} from 'lodash';
import {useState} from 'react';
import RunDetailAutomateDefinition from 'components/RunDetailAutomateDefinition';
import RunDetailAutomateMethods from 'components/RunDetailAutomateMethods';
import CliCommand from 'components/RunDetailAutomateMethods/methods/CLICommand';
import Cypress from 'components/RunDetailAutomateMethods/methods/Cypress';
import DeepLink from 'components/RunDetailAutomateMethods/methods/DeepLink';
import GithubActions from 'components/RunDetailAutomateMethods/methods/GithubActions';
import {CLI_RUNNING_TESTS_URL} from 'constants/Common.constants';
import {TriggerTypes} from 'constants/Test.constants';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import {useVariableSet} from 'providers/VariableSet';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunDetailAutomate.styled';

function getMethods(triggerType: TriggerTypes) {
  switch (triggerType) {
    case TriggerTypes.cypress:
      return [
        {
          id: 'cypress',
          label: 'Cypress',
          component: Cypress,
        },
      ];
    default:
      return [
        {
          id: 'cli',
          label: 'CLI',
          component: CliCommand,
        },
        {
          id: 'deeplink',
          label: 'Deep Link',
          component: DeepLink,
        },
        {
          id: 'githubAction',
          label: 'GitHub Actions',
          component: GithubActions,
        },
      ];
  }
}

interface IProps {
  test: Test;
  run: TestRun;
}

const RunDetailAutomate = ({test, run}: IProps) => {
  const [fileName, setFileName] = useState<string>(`${snakeCase(test.name)}.yaml`);
  const {selectedVariableSet: {id: variableSetId} = {}} = useVariableSet();

  return (
    <S.Container>
      <S.SectionLeft>
        <RunDetailAutomateDefinition
          id={test.id}
          version={test.version}
          resourceType={ResourceType.Test}
          fileName={fileName}
          onFileNameChange={setFileName}
        />
      </S.SectionLeft>
      <S.SectionRight>
        <RunDetailAutomateMethods
          resourceType={ResourceType.Test}
          methods={getMethods(test.trigger.type).map(({id, label, component: Component}) => ({
            id,
            label,
            children: (
              <Component
                docsUrl={CLI_RUNNING_TESTS_URL}
                fileName={fileName}
                id={test.id}
                resourceType={ResourceType.Test}
                run={run}
                test={test}
                variableSetId={variableSetId}
              />
            ),
          }))}
        />
      </S.SectionRight>
    </S.Container>
  );
};

export default RunDetailAutomate;
