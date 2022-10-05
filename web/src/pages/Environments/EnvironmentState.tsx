import {IEnvironment} from './IEnvironment';

export interface EnvironmentState {
  environment?: IEnvironment;
  query: string;
  dialog: boolean;
}
