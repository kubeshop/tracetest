import {EResourceType, TListResponse, TResource} from 'types/Settings.types';

export type TRawDemo = {
  id: string;
  name: string;
  otelCart: string;
  otelCheckout: string;
  otelEnabled: boolean;
  otelFrontend: string;
  otelProductCatalog: string;
  pokeshopEnabled: boolean;
  pokeshopGrpc: string;
  pokeshopHttp: string;
};

type Demo = {
  id: string;
  name: string;
  otelCart: string;
  otelCheckout: string;
  otelEnabled: boolean;
  otelFrontend: string;
  otelProductCatalog: string;
  pokeshopEnabled: boolean;
  pokeshopGrpc: string;
  pokeshopHttp: string;
};

function Demo(rawDemos?: TListResponse<TRawDemo>): Demo {
  const items = rawDemos?.items ?? [];
  const demo = items?.[0];

  return {
    id: demo?.spec?.id ?? '',
    name: demo?.spec?.name ?? '',
    otelCart: demo?.spec?.otelCart ?? '',
    otelCheckout: demo?.spec?.otelCheckout ?? '',
    otelEnabled: demo?.spec?.otelEnabled ?? false,
    otelFrontend: demo?.spec?.otelFrontend ?? '',
    otelProductCatalog: demo?.spec?.otelProductCatalog ?? '',
    pokeshopEnabled: demo?.spec?.pokeshopEnabled ?? false,
    pokeshopGrpc: demo?.spec?.pokeshopGrpc ?? '',
    pokeshopHttp: demo?.spec?.pokeshopHttp ?? '',
  };
}

export function rawToResource(spec: TRawDemo): TResource<TRawDemo> {
  return {spec, type: EResourceType.Demo};
}

export function isDemoEnabled(demo: Demo): boolean {
  return [demo.otelEnabled, demo.pokeshopEnabled].some(item => item);
}

export default Demo;
