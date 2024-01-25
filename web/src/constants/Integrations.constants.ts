export interface IIntegration {
  name: SupportedIntegrations;
  title: string;
  url: string;
  isActive: boolean;
  isAvailable: boolean;
}

export enum SupportedIntegrations {
  CYPRESS = 'Cypress',
  PLAYWRIGHT = 'Playwright',
  K6 = 'K6',
}

const Cypress: IIntegration = {
  name: SupportedIntegrations.CYPRESS,
  title: 'Cypress',
  url: 'https://docs.tracetest.io/tools-and-integrations/cypress',
  isActive: true,
  isAvailable: false,
};

const K6: IIntegration = {
  name: SupportedIntegrations.K6,
  title: 'K6',
  url: 'https://docs.tracetest.io/tools-and-integrations/k6',
  isActive: true,
  isAvailable: true,
};

const Playwright: IIntegration = {
  name: SupportedIntegrations.PLAYWRIGHT,
  title: 'Playwright',
  url: 'https://docs.tracetest.io/tools-and-integrations/playwright',
  isActive: true,
  isAvailable: false,
};

export const Integrations = {
  [SupportedIntegrations.CYPRESS]: Cypress,
  [SupportedIntegrations.PLAYWRIGHT]: Playwright,
  [SupportedIntegrations.K6]: K6,
} as const;
