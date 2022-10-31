import {TEnvironment, TRawEnvironment} from 'types/Environment.types';

function Environment({id = '', name = '', description = '', values = []}: TRawEnvironment): TEnvironment {
  return {
    id,
    name,
    description,
    values: values?.map(value => ({key: value?.key ?? '', value: value?.value ?? ''})),
  };
}

export default Environment;
