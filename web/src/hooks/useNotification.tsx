import {CheckCircleOutlined, CloseCircleOutlined, InfoCircleOutlined, WarningOutlined} from '@ant-design/icons';
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

const icons: Record<TNotificationType, React.ComponentType> = {
  info: InfoCircleOutlined,
  success: CheckCircleOutlined,
  error: CloseCircleOutlined,
  warning: WarningOutlined,
  open: InfoCircleOutlined,
};

const useNotification = () => {
  const [api, contextHolder] = notification.useNotification();
  const {notification: notificationStyles} = useTheme();

  const showNotification: TShowNotification = useCallback(
    ({type = 'info', title = '', ...rest}) => {
      const overwrite = notificationStyles[type];
      const notificationFn = api[type];
      const Icon = icons[type] as React.ComponentType<{style: React.CSSProperties}>;

      notificationFn({icon: <Icon style={{color: overwrite.color}} />, ...overwrite, ...rest, message: title});
    },
    [api, notificationStyles]
  );

  return {contextHolder, showNotification};
};

export default useNotification;
