import {TEnvironmentValue} from 'models/Environment.model';
import Test from 'models/Test.model';
import {getServerBaseUrl} from '../utils/Common';

export type TDeepLinkConfig = {
  variables: TEnvironmentValue[];
  useEnvironmentId: boolean;
  test: Test;
  environmentId?: string;
};

const DeepLinkService = () => ({
  getLink({variables, useEnvironmentId, test: {id: testId, version}, environmentId}: TDeepLinkConfig) {
    const baseUrl = getServerBaseUrl();
    const filteredVariables = variables.filter(variable => !!variable && variable.key);
    const stringVariables = encodeURI(JSON.stringify(filteredVariables));

    const url = `${baseUrl}test/${testId}/version/${version}/run?${
      filteredVariables.length ? `variables=${stringVariables}` : ''
    }${useEnvironmentId && environmentId ? `&environmentId=${environmentId}` : ''}`;

    return url;
  },
});

export default DeepLinkService();
