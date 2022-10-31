import {IKeyValue} from 'constants/Test.constants';
import {Model, TEnvironmentSchemas} from 'types/Common.types';

export type TRawEnvironment = TEnvironmentSchemas['Environment'];

export type TEnvironment = Model<TRawEnvironment, {values: IKeyValue[]}>;
