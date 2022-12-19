import {notification} from 'antd';
import {useCallback} from 'react';
import {useTheme} from 'styled-components';
import {SupportedDataStores, TConnectionResult} from 'types/Config.types';
import TestConnectionNotification from 'components/TestConnectionNotification/TestConnectionNotification';

const useTestConnectionNotification = () => {
  const [api, contextHolder] = notification.useNotification();
  const {
    notification: {success, error},
  } = useTheme();

  const showNotification = useCallback(
    (result: TConnectionResult, dataStoreType: SupportedDataStores) => {
      if (dataStoreType === SupportedDataStores.OtelCollector) {
        return api.info({
          message: 'No Automated Test',
          description:
            'Since the OpenTelemetry Collector sends traces to Tracetest, there is no automated test. Once you have configured your OpenTelemetry Collector to send Tracetest spans to Tracetest, try running a Tracetest test against your application under test.',
          style: {
            minWidth: success.style.minWidth,
          },
        });
      }

      if (result.allPassed) {
        return api.success({
          message: 'All tests successful - configuration is valid',
          description: <TestConnectionNotification result={result} />,
          ...success,
        });
      }

      api.error({
        message: 'Test failed - configuration is not valid',
        description: <TestConnectionNotification result={result} />,
        ...error,
      });
    },
    [api, error, success]
  );

  return {showNotification, contextHolder};
};

export default useTestConnectionNotification;
