import {TConnectionResult, TRawConnectionResult} from 'types/Config.types';
import ConnectionTestStep from './ConnectionResultStep.model';

const ConnectionResult = ({
  authentication: rawAuthentication = {},
  connectivity: rawConnectivity = {},
  fetchTraces: rawFetchTraces = {},
}: TRawConnectionResult): TConnectionResult => {
  const authentication = ConnectionTestStep(rawAuthentication);
  const connectivity = ConnectionTestStep(rawConnectivity);
  const fetchTraces = ConnectionTestStep(rawFetchTraces);

  return {
    allPassed: authentication.passed && connectivity.passed && fetchTraces.passed,
    authentication,
    connectivity,
    fetchTraces,
  };
};

export default ConnectionResult;
