import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import TracetestAPI from 'redux/apis/Tracetest';
import {TDraftResource} from 'types/Settings.types';
import {useNotification} from '../Notification/Notification.provider';

const {useCreateSettingMutation, useUpdateSettingMutation} = TracetestAPI.instance;

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
    async (resources: TDraftResource[]) => {
      await Promise.all(resources.map(resource => onSaveResource(resource)));

      showNotification({type: 'success', title: 'Settings saved', description: 'Your settings were saved'});
    },
    [onSaveResource, showNotification]
  );

  const value = useMemo<IContext>(
    () => ({isLoading: isLoadingCreate || isLoadingUpdate, onSubmit}),
    [isLoadingCreate, isLoadingUpdate, onSubmit]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SettingsProvider;
