import {useEffect, useState} from 'react';
import Validator from 'utils/Validator';

export type BodyMode = 'json' | 'xml' | 'raw' | 'none';

export function useBodyMode(body?: string) {
  const [bodyMode, setBodyMode] = useState<BodyMode>('none');
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    if (!initialized && body) {
      setInitialized(true);
      setBodyMode(Validator.getBodyType(body));
    }
  }, [setBodyMode, body, initialized, setInitialized]);

  return [bodyMode, setBodyMode] as const;
}
