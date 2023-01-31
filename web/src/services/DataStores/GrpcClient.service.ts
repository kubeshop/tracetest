import {IGRPCClientSettings, SupportedDataStores, TDataStoreService, TRawGRPCClientSettings} from 'types/Config.types';
import DataStore from 'models/DataStore.model';

const GrpcClientService = (): TDataStoreService => ({
  async getRequest({dataStore = {}}, dataStoreType = SupportedDataStores.JAEGER) {
    const values = dataStore[dataStoreType || SupportedDataStores.JAEGER] as IGRPCClientSettings;
    const {
      endpoint = '',
      readBufferSize,
      writeBufferSize,
      waitForReady = false,
      rawHeaders = [],
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
    const headers = rawHeaders.reduce((acc, curr) => ({...acc, [curr.key]: curr.value}), {});

    return Promise.resolve({
      type: dataStoreType,
      name: dataStoreType,
      [dataStoreType]: {
        endpoint,
        readBufferSize: parseInt(String(readBufferSize), 10),
        writeBufferSize: parseInt(String(writeBufferSize), 10),
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
  validateDraft({dataStore = {name: '', type: SupportedDataStores.JAEGER}, dataStoreType}) {
    const values = (dataStore[dataStoreType || SupportedDataStores.JAEGER] as IGRPCClientSettings) ?? {};
    const {endpoint = ''} = values;
    if (!endpoint) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues(
    {defaultDataStore = {name: '', type: SupportedDataStores.JAEGER} as DataStore},
    dataStoreType = SupportedDataStores.JAEGER
  ) {
    const values = (defaultDataStore[dataStoreType] as TRawGRPCClientSettings) ?? {};
    const {
      endpoint = '',
      readBufferSize,
      writeBufferSize,
      waitForReady = false,
      headers = {},
      balancerName = '',
      compression = '',
      tls: {
        insecure,
        insecureSkipVerify = false,
        serverName = '',
        settings: {cAFile = '', certFile = '', keyFile = '', minVersion = '', maxVersion = ''} = {},
      } = {},
      auth = {},
    } = values;

    const rawHeaders = Object.entries(headers).map(([key, value]) => ({key, value}));

    return {
      dataStore: {
        [dataStoreType]: {
          endpoint,
          readBufferSize,
          writeBufferSize,
          waitForReady,
          headers,
          rawHeaders,
          balancerName,
          compression,
          tls: {
            insecure: endpoint ? !!insecure : true,
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
        name: dataStoreType,
        type: dataStoreType,
      },
      dataStoreType,
    };
  },
});

export default GrpcClientService();
