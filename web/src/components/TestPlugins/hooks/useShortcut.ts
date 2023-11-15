import {Form} from 'antd';
import {FormInstance} from 'antd/es/form/Form';
import {useCallback, useEffect} from 'react';

export const useShortcutWithDefault = (form: FormInstance) => {
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

const useShortcut = () => {
  const form = Form.useFormInstance();

  useShortcutWithDefault(form);
};

export default useShortcut;
