import {IGRPCClientSettings, TRawGRPCClientSettings} from 'types/DataStore.types';

const GrpcClientService = () => ({
  async getRequest(values: IGRPCClientSettings): Promise<TRawGRPCClientSettings> {
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
    });
  },
  validateDraft(values: IGRPCClientSettings): Promise<boolean> {
    const {endpoint = ''} = values;
    if (!endpoint) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues(values: TRawGRPCClientSettings): IGRPCClientSettings {
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
    };
  },
});

export default GrpcClientService();
