import {Form, FormInstance} from 'antd';
import {useMemo} from 'react';
import {IValues} from '../TestSpecForm';

const useAssertionFormValues = (form: FormInstance<IValues>) => {
  const currentSelector = Form.useWatch('selector', form);
  const currentAssertions = Form.useWatch('assertions', form);

  return useMemo(
    () => ({
      currentSelector: currentSelector || '',
      currentAssertions: currentAssertions || [],
    }),
    [currentAssertions, currentSelector]
  );
};

export default useAssertionFormValues;
