import {Form} from 'antd';
import {useCallback, useEffect} from 'react';

const useShortcut = () => {
  const form = Form.useFormInstance();

  const onKeydown = useCallback(
    (e: KeyboardEvent) => {
      const modifierKey = e.metaKey || e.ctrlKey;
      if (e.key === 'Enter' && modifierKey) {
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
