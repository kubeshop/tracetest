import {SupportedDataStores, TDataStoreService} from 'types/DataStore.types';

const AwsXRayService = (): TDataStoreService => ({
  getRequest({
    dataStore: {
      awsxray: {region = '', accessKeyId = '', secretAccessKey = '', sessionToken = '', useDefaultAuth = false} = {},
    } = {},
  }) {
    return Promise.resolve({
      type: SupportedDataStores.AWSXRay,
      name: SupportedDataStores.AWSXRay,
      awsxray: {
        region,
        accessKeyId,
        secretAccessKey,
        sessionToken,
        useDefaultAuth,
      },
    });
  },
  validateDraft({
    dataStore: {awsxray: {region = '', accessKeyId = '', secretAccessKey = '', useDefaultAuth = false} = {}} = {},
  }) {
    if (((!accessKeyId || !secretAccessKey) && !useDefaultAuth) || !region) return Promise.resolve(false);

    return Promise.resolve(true);
  },
  getInitialValues({awsxray = {}}) {
    const {region = '', secretAccessKey = '', accessKeyId = '', sessionToken = '', useDefaultAuth = false} = awsxray;

    return {
      dataStore: {
        name: SupportedDataStores.AWSXRay,
        type: SupportedDataStores.AWSXRay,
        awsxray: {region, secretAccessKey, accessKeyId, sessionToken, useDefaultAuth},
      },
      dataStoreType: SupportedDataStores.AWSXRay,
    };
  },
  getIsOtlpBased() {
    return false;
  },
});

export default AwsXRayService();
