import {Model, TEnvironmentSchemas} from 'types/Common.types';
import {IKeyValue} from '../constants/Test.constants';

export type TRawEnvironment = TEnvironmentSchemas['Environment'];
export type TEnvironmentValue = TEnvironmentSchemas['EnvironmentValue'];
type Environment = Model<TRawEnvironment, {values: IKeyValue[]}>;

function Environment({id = '', name = '', description = '', values = []}: TRawEnvironment): Environment {
  return {
    id,
    name,
    description,
    values: values?.map(value => ({key: value?.key ?? '', value: value?.value ?? ''})),
  };
}

export default Environment;
