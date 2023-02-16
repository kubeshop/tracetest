import {SupportedDataStores, TDataStoreService} from 'types/Config.types';

const AwsXRayService = (): TDataStoreService => ({
  getRequest({
    dataStore: {awsxray: {region = '', accessKeyId = '', secretAccessKey = '', sessionToken = ''} = {}} = {},
  }) {
    return Promise.resolve({
      type: SupportedDataStores.AWSXRay,
      name: SupportedDataStores.AWSXRay,
      awsxray: {
        region,
        accessKeyId,
        secretAccessKey,
        sessionToken,
      },
    });
  },
  validateDraft({dataStore: {awsxray: {region = '', accessKeyId = '', secretAccessKey = ''} = {}} = {}}) {
    if (!region || !accessKeyId || !secretAccessKey) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({defaultDataStore: {awsxray = {}} = {}}) {
    const {region = '', secretAccessKey = '', accessKeyId = '', sessionToken = ''} = awsxray;

    return {
      dataStore: {
        name: SupportedDataStores.AWSXRay,
        type: SupportedDataStores.AWSXRay,
        awsxray: {region, secretAccessKey, accessKeyId, sessionToken},
      },
      dataStoreType: SupportedDataStores.AWSXRay,
    };
  },
});

export default AwsXRayService();
