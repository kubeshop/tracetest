import { Model, TConfigSchemas } from 'types/Common.types';
import ConnectionTestStep from './ConnectionResultStep.model';

export type TRawConnectionResult = TConfigSchemas['ConnectionResult'];
type ConnectionResult = Model<
  TRawConnectionResult,
  {
    allPassed: boolean;
    portCheck: ConnectionTestStep;
    authentication: ConnectionTestStep;
    connectivity: ConnectionTestStep;
    fetchTraces: ConnectionTestStep;
  }
>;

const ConnectionResult = ({
  portCheck: rawPortCheck = {},
  authentication: rawAuthentication = {},
  connectivity: rawConnectivity = {},
  fetchTraces: rawFetchTraces = {},
}: TRawConnectionResult): ConnectionResult => {
  const portCheck = ConnectionTestStep(rawPortCheck);
  const authentication = ConnectionTestStep(rawAuthentication);
  const connectivity = ConnectionTestStep(rawConnectivity);
  const fetchTraces = ConnectionTestStep(rawFetchTraces);

  return {
    allPassed: portCheck.status && authentication.passed && connectivity.passed && fetchTraces.passed,
    portCheck,
    authentication,
    connectivity,
    fetchTraces,
  };
};

export default ConnectionResult;
