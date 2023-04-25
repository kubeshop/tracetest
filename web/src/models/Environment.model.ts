import {Model, TEnvironmentSchemas} from 'types/Common.types';
import {IKeyValue} from '../constants/Test.constants';

export type TRawEnvironment = TEnvironmentSchemas['EnvironmentResource'];
export type TEnvironmentValue = TEnvironmentSchemas['EnvironmentValue'];
type Environment = Model<TEnvironmentSchemas['Environment'], {values: IKeyValue[]}>;

function Environment({spec: {id = '', name = '', description = '', values = []} = {}}: TRawEnvironment): Environment {
  return Environment.fromRun({id, name, description, values});
}

Environment.fromRun = ({
  id = '',
  name = '',
  description = '',
  values = [],
}: TEnvironmentSchemas['Environment']): Environment => {
  return {
    id,
    name,
    description,
    values: values?.map(value => ({key: value?.key ?? '', value: value?.value ?? ''})),
  };
};

export default Environment;
