import { Model, TConfigSchemas } from 'types/Common.types';
import ConnectionTestStep from './ConnectionResultStep.model';

export type TRawConnectionResult = TConfigSchemas['ConnectionResult'];
type ConnectionResult = Model<
  TRawConnectionResult,
  {
    allPassed: boolean;
    authentication: ConnectionTestStep;
    connectivity: ConnectionTestStep;
    fetchTraces: ConnectionTestStep;
  }
>;

const ConnectionResult = ({
  authentication: rawAuthentication = {},
  connectivity: rawConnectivity = {},
  fetchTraces: rawFetchTraces = {},
}: TRawConnectionResult): ConnectionResult => {
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
