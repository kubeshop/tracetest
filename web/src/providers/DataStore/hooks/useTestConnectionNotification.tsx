import {notification} from 'antd';
import {useCallback} from 'react';
import {useTheme} from 'styled-components';
import {TConnectionResult} from 'types/Config.types';
import TestConnectionNotification from 'components/TestConnectionNotification/TestConnectionNotification';

const useTestConnectionNotification = () => {
  const [api, contextHolder] = notification.useNotification();
  const {
    notification: {success, error},
  } = useTheme();

  const showNotification = useCallback(
    (result: TConnectionResult) => {
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
