import {useCallback} from 'react';
import {useNotification} from '../providers/Notification/Notification.provider';

const useCopy = () => {
  const {showNotification} = useNotification();
  const copy = useCallback(
    async (data: string) => {
      await navigator.clipboard.writeText(data);

      showNotification({
        type: 'success',
        title: 'Value copied to the clipboard',
      });
    },
    [showNotification]
  );

  return copy;
};

export default useCopy;
