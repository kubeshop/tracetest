import {EResourceType, TListResponse, TResource} from 'types/Settings.types';

export type TRawPolling = {
  default: boolean;
  id: string;
  name: string;
  periodic: {
    retryDelay: string;
    timeout: string;
  };
  strategy: string;
};

type Polling = {
  default: boolean;
  id: string;
  name: string;
  periodic: {
    retryDelay: string;
    timeout: string;
  };
  strategy: string;
};

function Polling(rawPollings?: TListResponse<TRawPolling>): Polling {
  const items = rawPollings?.items ?? [];
  const polling = items.find(item => item?.spec?.default ?? false);

  return {
    default: true,
    id: polling?.spec?.id ?? '',
    name: polling?.spec?.name ?? 'periodic',
    periodic: {
      retryDelay: polling?.spec?.periodic?.retryDelay ?? '',
      timeout: polling?.spec?.periodic?.timeout ?? '',
    },
    strategy: polling?.spec?.strategy ?? 'periodic',
  };
}

export function rawToResource(spec: TRawPolling): TResource<TRawPolling> {
  return {spec, type: EResourceType.Polling};
}

export default Polling;
