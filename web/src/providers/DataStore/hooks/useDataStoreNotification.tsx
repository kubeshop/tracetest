import {notification, Typography} from 'antd';
import {useCallback} from 'react';
import {useTheme} from 'styled-components';
import {SupportedDataStores, TConnectionResult} from 'types/DataStore.types';
import TestConnectionNotification from 'components/TestConnectionNotification/TestConnectionNotification';
import {NoTestConnectionDataStoreList} from '../../../constants/DataStore.constants';

const useDataStoreNotification = () => {
  const [api, contextHolder] = notification.useNotification();
  const {
    notification: {success, error},
  } = useTheme();

  const showTestConnectionNotification = useCallback(
    (result: TConnectionResult, dataStoreType: SupportedDataStores) => {
      if (NoTestConnectionDataStoreList.includes(dataStoreType)) {
        return api.info({
          message: <Typography.Title level={2}>No Automated Test</Typography.Title>,
          description:
            'Since the OpenTelemetry Collector sends traces to Tracetest, there is no automated test. Once you have configured your OpenTelemetry Collector to send Tracetest spans to Tracetest, try running a Tracetest test against your application under test.',
          style: {
            minWidth: success.style.minWidth,
          },
        });
      }

      if (result.allPassed) {
        return api.success({
          message: <Typography.Title level={2}>All tests successful - configuration is valid</Typography.Title>,
          description: <TestConnectionNotification result={result} />,
          ...success,
        });
      }

      api.error({
        message: <Typography.Title level={2}>Test failed - configuration is not valid</Typography.Title>,
        description: <TestConnectionNotification result={result} />,
        ...error,
      });
    },
    [api, error, success]
  );

  const showSuccessNotification = useCallback(() => {
    return api.success({
      message: 'Data Store saved',
      description: 'Your configuration was added',
      ...success,
    });
  }, [api, success]);

  return {contextHolder, showSuccessNotification, showTestConnectionNotification};
};

export default useDataStoreNotification;
