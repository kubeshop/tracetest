import {Dispatch, SetStateAction, useEffect} from 'react';

const key = 'tracetest-new-user';
const oldUserValue = String(true);

export function useOpenGuidedTourForNewUsersEffect(setVisible: Dispatch<SetStateAction<boolean>>): void {
  useEffect(() => {
    const isOldUser = localStorage.getItem(key);
    if (isOldUser !== oldUserValue) {
      setVisible(true);
      localStorage.setItem(key, oldUserValue);
    }
  }, [setVisible]);
}
