import {IGRPCClientSettings, SupportedDataStores, TDataStoreService, TRawGRPCClientSettings} from 'types/Config.types';

const GrpcClientService = (): TDataStoreService => ({
  async getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.JAEGER) {
    const values = dataStore[dataStoreType || SupportedDataStores.JAEGER] as IGRPCClientSettings;
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
        settings: {minVersion = '', maxVersion = ''} = {},
      } = {},
      auth = {},
      fileCA,
      fileCert,
      fileKey,
    } = values;

    const filesToText = [fileCA, fileCert, fileKey].map(file => (file ? file.text() : Promise.resolve(undefined)));
    const [cAFile, certFile, keyFile] = await Promise.all(filesToText);

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
    const values = dataStore[dataStoreType || SupportedDataStores.JAEGER] as IGRPCClientSettings;
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
          fileCA: cAFile ? new File([cAFile], 'fileCA') : undefined,
          fileCert: certFile ? new File([certFile], 'fileCert') : undefined,
          fileKey: keyFile ? new File([keyFile], 'fileKey') : undefined,
        },
      },
      dataStoreType,
    };
  },
});

export default GrpcClientService();
