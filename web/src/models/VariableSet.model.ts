import {Model, TVariableSetSchemas} from 'types/Common.types';
import {IKeyValue} from '../constants/Test.constants';

export type TRawVariableSet = TVariableSetSchemas['VariableSetResource'];
export type TVariableSetValue = TVariableSetSchemas['VariableSetValue'];
type VariableSet = Model<TVariableSetSchemas['VariableSet'], {values: IKeyValue[]}>;

function VariableSet({spec: {id = '', name = '', description = '', values = []} = {}}: TRawVariableSet): VariableSet {
  return VariableSet.fromRun({id, name, description, values});
}

VariableSet.fromRun = ({
  id = '',
  name = '',
  description = '',
  values = [],
}: TVariableSetSchemas['VariableSet']): VariableSet => {
  return {
    id,
    name,
    description,
    values: values?.map(value => ({key: value?.key ?? '', value: value?.value ?? ''})),
  };
};

export default VariableSet;
