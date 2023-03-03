import {IHttpClientSettings, TRawHttpClientSettings} from 'types/DataStore.types';

const HttpClientService = () => ({
  async getRequest(values: IHttpClientSettings): Promise<TRawHttpClientSettings> {
    const {
      url = '',
      rawHeaders = [],
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
    const headers = rawHeaders.reduce((acc, {key, value}) => (key && value ? {...acc, [key]: value} : acc), {});

    return {
      url,
      headers,
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
    };
  },
  validateDraft(values: IHttpClientSettings) {
    const {url = ''} = values;
    if (!url) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues(values: TRawHttpClientSettings): IHttpClientSettings {
    const {
      url = '',
      headers = {},
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
      url,
      headers,
      rawHeaders,
      tls: {
        insecure: url ? !!insecure : true,
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

export default HttpClientService();
