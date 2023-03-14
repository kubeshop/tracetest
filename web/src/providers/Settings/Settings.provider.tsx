import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import {useCreateSettingMutation, useUpdateSettingMutation} from 'redux/apis/TraceTest.api';
import {TDraftResource} from 'types/Settings.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useNotification} from '../Notification/Notification.provider';

interface IContext {
  isLoading: boolean;
  onSubmit(resource: TDraftResource[]): void;
}

const Context = createContext<IContext>({isLoading: false, onSubmit: noop});

interface IProps {
  children: React.ReactNode;
}

export const useSettings = () => useContext(Context);

const SettingsProvider = ({children}: IProps) => {
  const {showNotification} = useNotification();
  const [createSetting, {isLoading: isLoadingCreate}] = useCreateSettingMutation();
  const [updateSetting, {isLoading: isLoadingUpdate}] = useUpdateSettingMutation();
  const {onOpen: onOpenConfirmation} = useConfirmationModal();

  const onSaveResource = useCallback(
    async (resource: TDraftResource) => {
      if (resource.spec.id) {
        await updateSetting({resource});
      } else {
        await createSetting({resource});
      }
    },
    [createSetting, updateSetting]
  );

  const onSubmit = useCallback(
    (resources: TDraftResource[]) => {
      onOpenConfirmation({
        title: <p>Are you sure you want to save this Setting?</p>,
        heading: 'Save Confirmation',
        okText: 'Save',
        onConfirm: async () => {
          await Promise.all(resources.map(resource => onSaveResource(resource)));

          showNotification({type: 'success', title: 'Settings saved', description: 'Your settings were saved'});
        },
      });
    },
    [onOpenConfirmation, onSaveResource, showNotification]
  );

  const value = useMemo<IContext>(
    () => ({isLoading: isLoadingCreate || isLoadingUpdate, onSubmit}),
    [isLoadingCreate, isLoadingUpdate, onSubmit]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SettingsProvider;
