import {TOutput, TRawOutput} from 'types/Output.types';

const Output = ({id = '', selector = '', attribute = '', source = 'trigger', regex = ''}: TRawOutput): TOutput => {
  return {
    id,
    selector,
    attribute,
    source,
    regex,
  };
};

export default Output;
