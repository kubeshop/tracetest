import {Model, TTestSchemas} from './Common.types';

export type TRawSelectedSpans = TTestSchemas['SelectedSpansResult'];

export type TSelectedSpans = Model<TRawSelectedSpans, {}>;
