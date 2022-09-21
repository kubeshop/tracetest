import {Dispatch, SetStateAction, useEffect, useState} from 'react';
import Validator from 'utils/Validator';

export type BodyMode = 'json' | 'xml' | 'raw' | 'none';

function useGuessBodyModeEffect(setBodyMode: Dispatch<SetStateAction<BodyMode>>, body?: string) {
  const [initialized, setInitialized] = useState(false);
  useEffect(() => {
    if (!initialized && body) {
      setInitialized(true);
      setBodyMode(Validator.getBodyType(body));
    }
  }, [setBodyMode, body, initialized, setInitialized]);
}

export function useBodyMode(body?: string): [BodyMode, Dispatch<SetStateAction<BodyMode>>] {
  const [bodyMode, setBodyMode] = useState<BodyMode>('none');
  useGuessBodyModeEffect(setBodyMode, body);
  return [bodyMode, setBodyMode];
}
