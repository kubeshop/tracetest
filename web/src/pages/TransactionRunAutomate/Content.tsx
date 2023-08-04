import {snakeCase} from 'lodash';
import {useState} from 'react';
import RunDetailAutomateDefinition from 'components/RunDetailAutomateDefinition';
import RunDetailAutomateMethods from 'components/RunDetailAutomateMethods';
import CliCommand from 'components/RunDetailAutomateMethods/methods/CLICommand';
import {CLI_RUNNING_TRANSACTIONS_URL} from 'constants/Common.constants';
import useDocumentTitle from 'hooks/useDocumentTitle';
import {useVariableSet} from 'providers/VariableSet';
import {useTransaction} from 'providers/Transaction/Transaction.provider';
import {ResourceType} from 'types/Resource.type';
import * as S from './TransactionRunAutomate.styled';

const Content = () => {
  const {transaction} = useTransaction();
  useDocumentTitle(`${transaction.name} - Automate`);
  const [fileName, setFileName] = useState<string>(`${snakeCase(transaction.name)}.yaml`);
  const {selectedVariableSet: {id: variableSetId} = {}} = useVariableSet();

  return (
    <S.Container>
      <S.SectionLeft>
        <RunDetailAutomateDefinition
          id={transaction.id}
          version={transaction.version}
          resourceType={ResourceType.Transaction}
          fileName={fileName}
          onFileNameChange={setFileName}
        />
      </S.SectionLeft>
      <S.SectionRight>
        <RunDetailAutomateMethods
          resourceType={ResourceType.Transaction}
          methods={[
            {
              id: 'cli',
              label: 'CLI',
              children: (
                <CliCommand
                  id={transaction.id}
                  variableSetId={variableSetId}
                  fileName={fileName}
                  resourceType={ResourceType.Transaction}
                  docsUrl={CLI_RUNNING_TRANSACTIONS_URL}
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
