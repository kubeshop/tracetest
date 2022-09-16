import {Dispatch, SetStateAction, useEffect, useState} from 'react';
import Validator from 'utils/Validator';

export type BodyMode = 'json' | 'xml' | 'raw' | 'none';

function useGuessBodyModeEffect(
  isEditing: undefined | boolean,
  setBodyMode: Dispatch<SetStateAction<BodyMode>>,
  body?: string
) {
  const [initialized, setInitialized] = useState(false);
  useEffect(() => {
    if (!initialized && isEditing && body) {
      setInitialized(true);
      setBodyMode(Validator.getBodyType(body));
    }
  }, [isEditing, setBodyMode, body, initialized, setInitialized]);
}

export function useBodyMode(isEditing?: boolean, body?: string): [BodyMode, Dispatch<SetStateAction<BodyMode>>] {
  const [bodyMode, setBodyMode] = useState<BodyMode>('none');
  useGuessBodyModeEffect(isEditing, setBodyMode, body);
  return [bodyMode, setBodyMode];
}
