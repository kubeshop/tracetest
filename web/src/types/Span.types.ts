import {SemanticGroupNames} from '../constants/SemanticGroupNames.constants';
import {Model, TTraceSchemas} from './Common.types';
import {TSpanAttribute} from './SpanAttribute.types';

export type TSpanFlatAttribute = {
  key: string;
  value: string;
};

export type TRawSpan = TTraceSchemas['Span'];

export type TSpan = Model<
  TRawSpan,
  {
    attributes: Record<string, TSpanAttribute>;
    type: SemanticGroupNames;
    duration: number;
    signature: TSpanFlatAttribute[];
    attributeList: TSpanFlatAttribute[];
    children?: TSpan[];
  }
>;
