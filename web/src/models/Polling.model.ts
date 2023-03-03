import {EResourceType, TListResponse, TResource} from 'types/Settings.types';

export type TRawPolling = {
  id: string;
  maxWaitTimeForTrace: string;
  name: string;
  retryDelay: string;
};

type Polling = {
  id: string;
  maxWaitTimeForTrace: string;
  name: string;
  retryDelay: string;
};

function Polling(rawPollings?: TListResponse<TRawPolling>): Polling {
  const items = rawPollings?.items ?? [];
  const polling = items?.[0];

  return {
    id: polling?.spec?.id ?? '',
    maxWaitTimeForTrace: polling?.spec?.maxWaitTimeForTrace ?? '',
    name: polling?.spec?.name ?? '',
    retryDelay: polling?.spec?.retryDelay ?? '',
  };
}

export function rawToResource(spec: TRawPolling): TResource<TRawPolling> {
  return {spec, type: EResourceType.Polling};
}

export default Polling;
