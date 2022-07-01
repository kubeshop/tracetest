import {Form, FormInstance} from 'antd';
import {useMemo} from 'react';
import {IValues} from '../AssertionForm';

const useAssertionFormValues = (form: FormInstance<IValues>) => {
  const currentSelectorList = Form.useWatch('selectorList', form);
  const currentPseudoSelector = Form.useWatch('pseudoSelector', form);
  const currentIsAdvancedSelector = Form.useWatch('isAdvancedSelector', form);
  const currentSelector = Form.useWatch('selector', form);

  return useMemo(
    () => ({
      currentSelectorList: currentSelectorList || [],
      currentPseudoSelector: currentPseudoSelector || undefined,
      currentIsAdvancedSelector: currentIsAdvancedSelector || false,
      currentSelector: currentSelector || '',
    }),
    [currentIsAdvancedSelector, currentPseudoSelector, currentSelector, currentSelectorList]
  );
};

export default useAssertionFormValues;
