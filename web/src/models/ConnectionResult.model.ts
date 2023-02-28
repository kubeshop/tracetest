import { Model, TConfigSchemas } from 'types/Common.types';
import ConnectionTestStep from './ConnectionResultStep.model';

export type TRawConnectionResult = TConfigSchemas['ConnectionResult'];
type ConnectionResult = Model<
  TRawConnectionResult,
  {
    allPassed: boolean;
    endpointLinting: ConnectionTestStep;
    authentication: ConnectionTestStep;
    connectivity: ConnectionTestStep;
    fetchTraces: ConnectionTestStep;
  }
>;

const ConnectionResult = ({
  endpointLinting: rawEndpointLinting = {},
  authentication: rawAuthentication = {},
  connectivity: rawConnectivity = {},
  fetchTraces: rawFetchTraces = {},
}: TRawConnectionResult): ConnectionResult => {
  const endpointLinting = ConnectionTestStep(rawEndpointLinting);
  const authentication = ConnectionTestStep(rawAuthentication);
  const connectivity = ConnectionTestStep(rawConnectivity);
  const fetchTraces = ConnectionTestStep(rawFetchTraces);

  return {
    allPassed: endpointLinting.status && authentication.passed && connectivity.passed && fetchTraces.passed,
    endpointLinting,
    authentication,
    connectivity,
    fetchTraces,
  };
};

export default ConnectionResult;
