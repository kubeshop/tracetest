import {IEnvironment} from '../../redux/apis/TraceTest.api';

export interface EnvironmentState {
  environment?: IEnvironment;
  query: string;
  dialog: boolean;
}
