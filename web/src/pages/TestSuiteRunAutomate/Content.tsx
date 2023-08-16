import {snakeCase} from 'lodash';
import {useState} from 'react';
import RunDetailAutomateDefinition from 'components/RunDetailAutomateDefinition';
import RunDetailAutomateMethods from 'components/RunDetailAutomateMethods';
import CliCommand from 'components/RunDetailAutomateMethods/methods/CLICommand';
import {CLI_RUNNING_TEST_SUITES_URL} from 'constants/Common.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import {useVariableSet} from 'providers/VariableSet';
import {useTestSuite} from 'providers/TestSuite';
import {ResourceType} from 'types/Resource.type';
import * as S from './TestSuiteRunAutomate.styled';

const Content = () => {
  const {testSuite} = useTestSuite();
  useDocumentTitle(`${testSuite.name} - Automate`);
  const [fileName, setFileName] = useState<string>(`${snakeCase(testSuite.name)}.yaml`);
  const {selectedVariableSet: {id: variableSetId} = {}} = useVariableSet();

  return (
    <S.Container>
      <S.SectionLeft>
        <RunDetailAutomateDefinition
          id={testSuite.id}
          version={testSuite.version}
          resourceType={ResourceType.TestSuite}
          fileName={fileName}
          onFileNameChange={setFileName}
        />
      </S.SectionLeft>
      <S.SectionRight>
        <RunDetailAutomateMethods
          resourceType={ResourceType.TestSuite}
          methods={[
            {
              id: 'cli',
              label: 'CLI',
              children: (
                <CliCommand
                  id={testSuite.id}
                  variableSetId={variableSetId}
                  fileName={fileName}
                  resourceType={ResourceType.TestSuite}
                  docsUrl={CLI_RUNNING_TEST_SUITES_URL}
                />
              ),
            },
          ]}
        />
      </S.SectionRight>
    </S.Container>
  );
};

export default Content;
