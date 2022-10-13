import {Model} from './Common.types';

export type TRawOutput = {
  id?: string;
  source?: string;
  attribute?: string;
  regex?: string;
  selector?: string;
};

export type TOutput = Model<TRawOutput, {
  regex?: string;
  selector?: string;
}>;
