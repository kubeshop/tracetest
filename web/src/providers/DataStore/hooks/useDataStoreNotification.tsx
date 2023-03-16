import {Typography} from 'antd';
import {useCallback} from 'react';
import {SupportedDataStores, TConnectionResult} from 'types/DataStore.types';
import TestConnectionNotification from 'components/TestConnectionNotification/TestConnectionNotification';
import {NoTestConnectionDataStoreList} from 'constants/DataStore.constants';
import {useNotification} from 'providers/Notification/Notification.provider';

const useDataStoreNotification = () => {
  const {showNotification} = useNotification();

  const showTestConnectionNotification = useCallback(
    (result: TConnectionResult, dataStoreType: SupportedDataStores) => {
      if (NoTestConnectionDataStoreList.includes(dataStoreType)) {
        return showNotification({
          type: 'info',
          title: <Typography.Title level={2}>No Automated Test</Typography.Title>,
          description:
            'Since the OpenTelemetry Collector sends traces to Tracetest, there is no automated test. Once you have configured your OpenTelemetry Collector to send Tracetest spans to Tracetest, try running a Tracetest test against your application under test.',
        });
      }

      if (result.allPassed) {
        return showNotification({
          type: 'success',
          title: <Typography.Title level={2}>All tests successful - configuration is valid</Typography.Title>,
          description: <TestConnectionNotification result={result} />,
        });
      }

      showNotification({
        type: 'error',
        title: <Typography.Title level={2}>Test failed - configuration is not valid</Typography.Title>,
        description: <TestConnectionNotification result={result} />,
      });
    },
    [showNotification]
  );

  const showSuccessNotification = useCallback(() => {
    return showNotification({
      type: 'success',
      title: 'Data Store saved',
      description: 'Your configuration was added',
    });
  }, [showNotification]);

  return {showSuccessNotification, showTestConnectionNotification};
};

export default useDataStoreNotification;
