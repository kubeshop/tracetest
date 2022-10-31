import {uniq} from 'lodash';
import {TRawSelectedSpans, TSelectedSpans} from 'types/SelectedSpans.types';

function SelectedSpans(rawSelectedSpans: TRawSelectedSpans): TSelectedSpans {
  return {
    selector: rawSelectedSpans?.selector ?? {},
    spanIds: uniq(rawSelectedSpans?.spanIds ?? [])
  };
}

export default SelectedSpans;
