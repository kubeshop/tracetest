import {noop} from 'lodash';
import {createContext, useContext, useMemo} from 'react';
import useNotificationHook, {TShowNotification} from 'hooks/useNotification';

interface IContext {
  showNotification: TShowNotification;
}

export const Context = createContext<IContext>({
  showNotification: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useNotification = () => useContext(Context);

const NotificationProvider = ({children}: IProps) => {
  const {contextHolder, showNotification} = useNotificationHook();

  const value = useMemo<IContext>(
    () => ({
      showNotification,
    }),
    [showNotification]
  );

  return (
    <>
      {contextHolder}
      <Context.Provider value={value}>{children}</Context.Provider>
    </>
  );
};

export default NotificationProvider;
