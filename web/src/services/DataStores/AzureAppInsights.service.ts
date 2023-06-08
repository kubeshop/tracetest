import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const AzureAppInsightsService = (): TDataStoreService => ({
  getRequest({dataStore: {azureappinsights: {resourceArmId = ''} = {}} = {}}) {
    return Promise.resolve({
      type: SupportedDataStores.AzureAppInsights,
      name: SupportedDataStores.AzureAppInsights,
      azureappinsights: {
        resourceArmId,
      },
    });
  },
  validateDraft({dataStore: {azureappinsights: {resourceArmId = ''} = {}} = {}}) {
    if (!resourceArmId) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({defaultDataStore: {azureappinsights = {}} = {}}) {
    const {resourceArmId = ''} = azureappinsights;

    return {
      dataStore: {
        name: SupportedDataStores.AzureAppInsights,
        type: SupportedDataStores.AzureAppInsights,
        azureappinsights: {resourceArmId},
      },
      dataStoreType: SupportedDataStores.AzureAppInsights,
    };
  },
});

export default AzureAppInsightsService();
