import {Form, FormInstance} from 'antd';
import {useMemo} from 'react';
import {IValues} from '../TestSpecForm';

const useAssertionFormValues = (form: FormInstance<IValues>) => {
  const currentSelector = Form.useWatch('selector', form);
  const currentAssertionList = Form.useWatch('assertionList', form);

  return useMemo(
    () => ({
      currentSelector: currentSelector || '',
      currentAssertionList: currentAssertionList || [],
    }),
    [currentAssertionList, currentSelector]
  );
};

export default useAssertionFormValues;
