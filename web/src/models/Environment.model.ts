import {TEnvironment, TRawEnvironment} from 'types/Environment.types';

const Environment = ({id = '', name = '', description = '', variables = []}: TRawEnvironment): TEnvironment => ({
  id,
  name,
  description,
  variables,
});

export default Environment;
