export interface IDraftSettings {
  analytics: boolean;
  maxWaitTimeForTrace: string;
  retryDelay: string;
  demo: IDemoSettings;
}

interface IDemoSettings {
  pokeshopEnabled: boolean;
  pokeshopHttp: string;
  pokeshopGrpc: string;
  otelEnabled: boolean;
  otelFrontend: string;
  otelProductCatalog: string;
  otelCart: string;
  otelCheckout: string;
}
