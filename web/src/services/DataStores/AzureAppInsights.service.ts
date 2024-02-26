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

    return Promise.resolve(true);
  },
  getInitialValues({azureappinsights = {}}) {
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
        azureappinsights: {
          resourceArmId,
          connectionType,
          accessToken,
          useAzureActiveDirectoryAuth,
        },
      },
      dataStoreType: SupportedDataStores.AzureAppInsights,
    };
  },
  getIsOtlpBased({dataStore: {azureappinsights: {connectionType = ConnectionTypes.Direct} = {}} = {}}) {
    return connectionType === ConnectionTypes.Collector;
  },
  getPublicInfo({azureappinsights = {}}) {
    const {
      resourceArmId = '',
      connectionType = ConnectionTypes.Direct,
      useAzureActiveDirectoryAuth = true,
    } = azureappinsights;

    return {
      'Resource Arm Id': resourceArmId,
      'Connection Type': connectionType,
      'Use Azure Active Directory Auth': useAzureActiveDirectoryAuth ? 'Yes' : 'No',
    };
  },
});

export default AzureAppInsightsService();
