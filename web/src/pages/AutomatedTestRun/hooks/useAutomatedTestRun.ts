import {useCallback} from 'react';
import {TVariableSetValue} from 'models/VariableSet.model';
import {useTest} from 'providers/Test/Test.provider';
import {useVariableSet} from 'providers/VariableSet';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';

const getParsedVariables = (rawVars: string): TVariableSetValue[] => {
  try {
    const variables = JSON.parse(rawVars) as TVariableSetValue[];

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
  const {selectedVariableSet} = useVariableSet();
  const {navigate} = useDashboard();

  const onAutomatedRun = useCallback(() => {
    const variables = getParsedVariables(query.get('variables') ?? '[]');
    const variableSetId = query.get('variableSetId') ?? selectedVariableSet?.id;

    onRun({
      variables,
      variableSetId,
      onCancel() {
        navigate(`/test/${testId}`, {replace: true});
      },
    });
  }, [navigate, onRun, query, selectedVariableSet?.id, testId]);

  return onAutomatedRun;
};

export default useAutomatedTestRun;
