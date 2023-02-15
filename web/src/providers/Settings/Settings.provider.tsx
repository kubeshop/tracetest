import {notification} from 'antd';
import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {useTheme} from 'styled-components';

import {IDraftSettings} from 'types/Settings.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  onSubmit(values: IDraftSettings): void;
}

const Context = createContext<IContext>({onSubmit: noop});

interface IProps {
  children: React.ReactNode;
}

export const useSettings = () => useContext(Context);

const SettingsProvider = ({children}: IProps) => {
  const [notificationApi, notificationComponent] = notification.useNotification();
  const {onOpen: onOpenConfirmation} = useConfirmationModal();
  const {
    notification: {success},
  } = useTheme();

  const onSubmit = useCallback(
    (values: IDraftSettings) => {
      onOpenConfirmation({
        title: <p>Are you sure you want to save this Setting?</p>,
        heading: 'Save Confirmation',
        okText: 'Save',
        onConfirm: () => {
          notificationApi.success({
            message: 'Settings saved',
            description: 'Your settings were saved',
            ...success,
          });
        },
      });
    },
    [notificationApi, onOpenConfirmation, success]
  );

  const value = useMemo<IContext>(() => ({onSubmit}), [onSubmit]);

  return (
    <>
      {notificationComponent}
      <Context.Provider value={value}>{children}</Context.Provider>
    </>
  );
};

export default SettingsProvider;
