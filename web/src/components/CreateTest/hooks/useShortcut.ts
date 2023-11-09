import {Form} from 'antd';
import {useCallback, useEffect} from 'react';

const useShortcut = () => {
  const form = Form.useFormInstance();

  const onKeydown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === 'Enter' && e.metaKey) {
        form.submit();
      }
    },
    [form]
  );

  useEffect(() => {
    window.addEventListener('keydown', onKeydown);
    return () => {
      window.removeEventListener('keydown', onKeydown);
    };
  }, [onKeydown]);
};

export default useShortcut;
