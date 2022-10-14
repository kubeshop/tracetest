import {IKeyValue} from 'constants/Test.constants';
import {Model} from 'types/Common.types';

export type TRawEnvironment = {
  id?: string;
  name?: string;
  description?: string;
  variables?: IKeyValue[];
}

export type TEnvironment = Model<TRawEnvironment, {}>;
