import {useCallback} from 'react';
import {TEnvironmentValue} from 'models/Environment.model';
import {useTest} from 'providers/Test/Test.provider';
import {useEnvironment} from 'providers/Environment/Environment.provider';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

const getParsedVariables = (rawVars: string): TEnvironmentValue[] => {
  try {
    const variables = JSON.parse(rawVars) as TEnvironmentValue[];

    return Array.isArray(variables) ? variables : [];
  } catch (err) {
    return [];
  }
};

const useAutomatedTestRun = (query: URLSearchParams) => {
  const {
    onRun,
    test: {id: testId},
  } = useTest();
  const {selectedEnvironment} = useEnvironment();
  const {navigate} = useDashboard();

  const onAutomatedRun = useCallback(() => {
    const variables = getParsedVariables(query.get('variables') ?? '[]');
    const environmentId = query.get('environmentId') ?? selectedEnvironment?.id;

    onRun({
      variables,
      environmentId,
      onCancel() {
        navigate(`/test/${testId}`, {replace: true});
      },
    });
  }, [navigate, onRun, query, selectedEnvironment?.id, testId]);

  return onAutomatedRun;
};

export default useAutomatedTestRun;
