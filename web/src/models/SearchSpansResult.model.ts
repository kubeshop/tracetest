import {Model, TTestSchemas} from '../types/Common.types';

export type TRawSearchSpansResult = TTestSchemas['SearchSpansResult'];
type SearchSpansResult = Model<
  TRawSearchSpansResult,
  {
    spanIds: string[];
    spansIds?: undefined;
  }
>;

const defaultSearchSpansResult: TRawSearchSpansResult = {
  spansIds: [],
};

function SearchSpansResult({spansIds = []} = defaultSearchSpansResult): SearchSpansResult {
  return {
    spanIds: spansIds,
  };
}

export default SearchSpansResult;
