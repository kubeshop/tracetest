import {TSpanFlatAttribute} from '../types/Span.types';
import {TSpanAttribute} from '../types/SpanAttribute.types';

const SpanAttribute = ({key, value}: TSpanFlatAttribute): TSpanAttribute => {
  return {
    name: key,
    value,
  };
};

export default SpanAttribute;
