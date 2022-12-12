import {SupportedDataStores, TDataStoreService, TRawGRPCClientSettings} from 'types/Config.types';

const GrpcClientService = (): TDataStoreService => ({
  getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.JAEGER) {
    const values = dataStore[dataStoreType || SupportedDataStores.JAEGER] as TRawGRPCClientSettings;
    const {
      endpoint = '',
      readBufferSize,
      writeBufferSize,
      waitForReady = false,
      headers = [],
      balancerName = '',
      compression = '',
      tls: {
        insecure = true,
        insecureSkipVerify = false,
        serverName = '',
        settings: {cAFile = '', certFile = '', keyFile = '', minVersion = '', maxVersion = ''} = {},
      } = {},
      auth = {},
    } = values;

    return Promise.resolve({
      type: dataStoreType,
      [dataStoreType]: {
        endpoint,
        readBufferSize,
        writeBufferSize,
        waitForReady,
        headers,
        balancerName,
        compression,
        tls: {
          insecure,
          insecureSkipVerify,
          serverName,
          settings: {
            cAFile,
            certFile,
            keyFile,
            minVersion,
            maxVersion,
          },
        },
        auth,
      },
    });
  },
  validateDraft({dataStore = {}, dataStoreType}) {
    const values = dataStore[dataStoreType || SupportedDataStores.JAEGER] as TRawGRPCClientSettings;
    const {endpoint = ''} = values;
    if (!endpoint) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({dataStores = []}, dataStoreType = SupportedDataStores.JAEGER) {
    const [dataStore = {}] = dataStores;
    const values = (dataStore[dataStoreType] as TRawGRPCClientSettings) ?? {};
    const {
      endpoint = '',
      readBufferSize,
      writeBufferSize,
      waitForReady = false,
      headers = [],
      balancerName = '',
      compression = '',
      tls: {
        insecure = true,
        insecureSkipVerify = false,
        serverName = '',
        settings: {cAFile = '', certFile = '', keyFile = '', minVersion = '', maxVersion = ''} = {},
      } = {},
      auth = {},
    } = values;

    return {
      dataStore: {
        [dataStoreType]: {
          endpoint,
          readBufferSize,
          writeBufferSize,
          waitForReady,
          headers,
          balancerName,
          compression,
          tls: {
            insecure,
            insecureSkipVerify,
            serverName,
            settings: {
              cAFile,
              certFile,
              keyFile,
              minVersion,
              maxVersion,
            },
          },
          auth,
        },
      },
      dataStoreType,
    };
  },
});

export default GrpcClientService();
