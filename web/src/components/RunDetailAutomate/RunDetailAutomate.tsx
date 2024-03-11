import {snakeCase} from 'lodash';
import {useEffect, useState} from 'react';
import RunDetailAutomateDefinition from 'components/RunDetailAutomateDefinition';
import RunDetailAutomateMethods from 'components/RunDetailAutomateMethods';
import CliCommand from 'components/RunDetailAutomateMethods/methods/CLICommand';
import Cypress from 'components/RunDetailAutomateMethods/methods/Cypress';
import DeepLink from 'components/RunDetailAutomateMethods/methods/DeepLink';
import Playwright from 'components/RunDetailAutomateMethods/methods/Playwright';
import Typescript from 'components/RunDetailAutomateMethods/methods/Typescript';
import GithubActions from 'components/RunDetailAutomateMethods/methods/GithubActions';
import K6 from 'components/RunDetailAutomateMethods/methods/K6';
import Artillery from 'components/RunDetailAutomateMethods/methods/Artillery';
import ArtilleryEngine from 'components/RunDetailAutomateMethods/methods/ArtilleryEngine';
import {CLI_RUNNING_TESTS_URL} from 'constants/Common.constants';
import useDefinitionFile from 'hooks/useDefinitionFile';
import {TriggerTypes} from 'constants/Test.constants';
import Test from 'models/Test.model';
import TestRun from 'models/TestRun.model';
import {useVariableSet} from 'providers/VariableSet';
import {ResourceType} from 'types/Resource.type';
import * as S from './RunDetailAutomate.styled';
import ResizablePanels, {FillPanel, RightPanel} from '../ResizablePanels';

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
    case TriggerTypes.playwright:
      return [
        {
          id: 'playwright',
          label: 'Playwright',
          component: Playwright,
        },
      ];
    case TriggerTypes.artillery:
      return [
        {
          id: 'artillery',
          label: 'Artillery',
          component: Artillery,
        },
      ];
    case TriggerTypes.k6:
      return [
        {
          id: 'k6',
          label: 'K6',
          component: K6,
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
        {
          id: 'typescript',
          label: 'TypeScript',
          component: Typescript,
        },
        {
          id: 'artillery-engine',
          label: 'Artillery',
          component: ArtilleryEngine,
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
  const {definition, loadDefinition} = useDefinitionFile();
  const {selectedVariableSet: {id: variableSetId} = {}} = useVariableSet();

  useEffect(() => {
    loadDefinition(ResourceType.Test, test.id, test.version);
  }, [loadDefinition, test.id, test.version]);

  return (
    <S.Container>
      <ResizablePanels saveId="run-detail-automate">
        <FillPanel>
          <RunDetailAutomateDefinition
            definition={definition}
            resourceType={ResourceType.Test}
            fileName={fileName}
            onFileNameChange={setFileName}
          />
        </FillPanel>
        <RightPanel
          panel={{
            openSize: () => (window.innerWidth / 2 / window.innerWidth) * 100,
            isDefaultOpen: true,
          }}
        >
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
                  definition={definition}
                  variableSetId={variableSetId}
                />
              ),
            }))}
          />
        </RightPanel>
      </ResizablePanels>
    </S.Container>
  );
};

export default RunDetailAutomate;
