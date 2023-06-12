import {ConnectionTypes, SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const AzureAppInsightsService = (): TDataStoreService => ({
  getRequest({
    dataStore: {
      azureappinsights: {
        resourceArmId = '',
        connectionType = ConnectionTypes.Direct,
        useAzureActiveDirectoryAuth = true,
        accessToken = '',
      } = {},
    } = {},
  }) {
    return Promise.resolve({
      type: SupportedDataStores.AzureAppInsights,
      name: SupportedDataStores.AzureAppInsights,
      azureappinsights: {
        resourceArmId,
        connectionType,
        useAzureActiveDirectoryAuth,
        accessToken,
      },
    });
  },
  validateDraft({
    dataStore: {
      isIngestorEnabled = false,
      azureappinsights: {
        resourceArmId = '',
        connectionType = ConnectionTypes.Direct,
        accessToken = '',
        useAzureActiveDirectoryAuth = true,
      } = {},
    } = {},
  }) {
    if (connectionType === ConnectionTypes.Direct && !resourceArmId) return Promise.resolve(false);
    if (connectionType === ConnectionTypes.Direct && !useAzureActiveDirectoryAuth && !accessToken)
      return Promise.resolve(false);
    if (connectionType === ConnectionTypes.Collector && !isIngestorEnabled) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({defaultDataStore: {azureappinsights = {}} = {}}, dataStoreType, configuredDataStore) {
    const {
      resourceArmId = '',
      connectionType = ConnectionTypes.Direct,
      accessToken = '',
      useAzureActiveDirectoryAuth = true,
    } = azureappinsights;

    return {
      dataStore: {
        name: SupportedDataStores.AzureAppInsights,
        type: SupportedDataStores.AzureAppInsights,
        azureappinsights: {resourceArmId, connectionType, accessToken, useAzureActiveDirectoryAuth},
        isIngestorEnabled:
          configuredDataStore === SupportedDataStores.AzureAppInsights && connectionType === ConnectionTypes.Collector,
      },
      dataStoreType: SupportedDataStores.AzureAppInsights,
    };
  },
  shouldTestConnection({dataStore: {azureappinsights: {connectionType = ConnectionTypes.Direct} = {}} = {}}) {
    return connectionType === ConnectionTypes.Direct;
  },
});

export default AzureAppInsightsService();
