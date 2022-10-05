import {IKeyValue} from 'constants/Test.constants';

export interface IEnvironment {
  id: string;
  name: string;
  description: string;
  variables?: IKeyValue[];
}
