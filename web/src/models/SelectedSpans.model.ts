import {uniq} from 'lodash';
import {Model, TTestSchemas} from 'types/Common.types';

export type TRawSelectedSpans = TTestSchemas['SelectedSpansResult'];
type SelectedSpans = Model<TRawSelectedSpans, {}>;

function SelectedSpans(rawSelectedSpans: TRawSelectedSpans): SelectedSpans {
  return {
    selector: rawSelectedSpans?.selector ?? {},
    spanIds: uniq(rawSelectedSpans?.spanIds ?? []),
  };
}

export default SelectedSpans;
