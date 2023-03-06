import {notification} from 'antd';
import {NotificationInstance, ArgsProps} from 'antd/lib/notification/index';
import {useCallback} from 'react';
import {useTheme} from 'styled-components';

export type TNotificationType = keyof NotificationInstance;
export type IShowNotificationProps = Omit<ArgsProps, 'message' | 'type'> & {
  title: React.ReactNode;
  type: TNotificationType;
};
export type TShowNotification = (props: IShowNotificationProps) => void;

const useNotification = () => {
  const [api, contextHolder] = notification.useNotification();
  const {notification: notificationStyles} = useTheme();

  const showNotification: TShowNotification = useCallback(
    ({type = 'info', title = '', ...rest}) => {
      const overwrite = notificationStyles[type];
      const notificationFn = api[type];

      notificationFn({...overwrite, ...rest, message: title});
    },
    [api, notificationStyles]
  );

  return {contextHolder, showNotification};
};

export default useNotification;
