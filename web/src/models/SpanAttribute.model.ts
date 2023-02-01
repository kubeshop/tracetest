import {TSpanFlatAttribute} from 'types/Span.types';

type SpanAttribute = {
  name: string;
  value: string;
};

const SpanAttribute = ({key, value}: TSpanFlatAttribute): SpanAttribute => {
  return {
    name: key,
    value,
  };
};

export default SpanAttribute;
