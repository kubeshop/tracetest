import {TVariableSetValue} from 'models/VariableSet.model';
import Test from 'models/Test.model';
import {getServerBaseUrl} from 'utils/Common';

export type TDeepLinkConfig = {
  variables: TVariableSetValue[];
  useVariableSetId: boolean;
  test: Test;
  variableSetId?: string;
  baseUrl?: string;
};

const DeepLinkService = () => ({
  getLink({
    baseUrl = getServerBaseUrl(),
    variables,
    useVariableSetId,
    test: {id: testId},
    variableSetId,
  }: TDeepLinkConfig) {
    const filteredVariables = variables.filter(variable => !!variable && variable.key);
    const stringVariables = encodeURI(JSON.stringify(filteredVariables));

    const url = `${baseUrl}test/${testId}/run?${filteredVariables.length ? `variables=${stringVariables}` : ''}${
      useVariableSetId && variableSetId ? `&variableSetId=${variableSetId}` : ''
    }`;

    return url;
  },
});

export default DeepLinkService();
