import {Dispatch, SetStateAction, useEffect} from 'react';

const key = 'tracetest-new-user';
const value = String(true);

export function useOpenGuidedTourForNewUsersEffect(setVisible: Dispatch<SetStateAction<boolean>>): void {
  useEffect(() => {
    const isNewUser = localStorage.getItem(key);
    if (isNewUser !== value) {
      setVisible(true);
      localStorage.setItem(key, value);
    }
  }, [setVisible]);
}
