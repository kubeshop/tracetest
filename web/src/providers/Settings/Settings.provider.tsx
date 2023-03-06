import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import {IDraftSettings} from 'types/Settings.types';
import {useConfirmationModal} from '../ConfirmationModal/ConfirmationModal.provider';
import {useNotification} from '../Notification/Notification.provider';

interface IContext {
  onSubmit(values: IDraftSettings): void;
}

const Context = createContext<IContext>({onSubmit: noop});

interface IProps {
  children: React.ReactNode;
}

export const useSettings = () => useContext(Context);

const SettingsProvider = ({children}: IProps) => {
  const {showNotification} = useNotification();
  const {onOpen: onOpenConfirmation} = useConfirmationModal();

  const onSubmit = useCallback(
    (values: IDraftSettings) => {
      onOpenConfirmation({
        title: <p>Are you sure you want to save this Setting?</p>,
        heading: 'Save Confirmation',
        okText: 'Save',
        onConfirm: () => {
          showNotification({type: 'success', title: 'Settings saved', description: 'Your settings were saved'});
        },
      });
    },
    [onOpenConfirmation, showNotification]
  );

  const value = useMemo<IContext>(() => ({onSubmit}), [onSubmit]);

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SettingsProvider;
