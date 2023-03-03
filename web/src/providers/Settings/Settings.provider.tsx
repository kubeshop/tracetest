import {notification} from 'antd';
import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {useTheme} from 'styled-components';

import {useCreateSettingMutation, useUpdateSettingMutation} from 'redux/apis/TraceTest.api';
import {TResource, TSpec} from 'types/Settings.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';

interface IContext {
  isLoading: boolean;
  onSubmit(resource: TResource<TSpec>): void;
}

const Context = createContext<IContext>({isLoading: false, onSubmit: noop});

interface IProps {
  children: React.ReactNode;
}

export const useSettings = () => useContext(Context);

const SettingsProvider = ({children}: IProps) => {
  const [createSetting, {isLoading: isLoadingCreate}] = useCreateSettingMutation();
  const [updateSetting, {isLoading: isLoadingUpdate}] = useUpdateSettingMutation();
  const [notificationApi, notificationComponent] = notification.useNotification();
  const {onOpen: onOpenConfirmation} = useConfirmationModal();
  const {
    notification: {success},
  } = useTheme();

  const onSubmit = useCallback(
    (resource: TResource<TSpec>) => {
      onOpenConfirmation({
        title: <p>Are you sure you want to save this Setting?</p>,
        heading: 'Save Confirmation',
        okText: 'Save',
        onConfirm: async () => {
          if (resource.spec.id) {
            await updateSetting({resource});
          } else {
            await createSetting({resource});
          }

          notificationApi.success({
            message: 'Settings saved',
            description: 'Your settings were saved',
            ...success,
          });
        },
      });
    },
    [createSetting, notificationApi, onOpenConfirmation, success, updateSetting]
  );

  const value = useMemo<IContext>(
    () => ({isLoading: isLoadingCreate || isLoadingUpdate, onSubmit}),
    [isLoadingCreate, isLoadingUpdate, onSubmit]
  );

  return (
    <>
      {notificationComponent}
      <Context.Provider value={value}>{children}</Context.Provider>
    </>
  );
};

export default SettingsProvider;
