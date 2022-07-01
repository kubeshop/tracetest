import {Form, FormInstance} from 'antd';
import {FieldData} from 'antd/node_modules/rc-field-form/es/interface';
import {isEmpty} from 'lodash';
import {useCallback} from 'react';
import SelectorService from 'services/Selector.service';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {TAssertion} from 'types/Assertion.types';
import {IValues} from '../AssertionForm';
import {TSpanFlatAttribute} from '../../../types/Span.types';

const {onChecksChange, onSelectorChange} = CreateAssertionModalAnalyticsService;

interface IProps {
  form: FormInstance<IValues>;
  attributeList: TSpanFlatAttribute[];
}

const useOnFieldsChange = ({form, attributeList}: IProps) => {
  const currentPseudoSelector = Form.useWatch('pseudoSelector', form) || undefined;
  const currentSelectorList = Form.useWatch('selectorList', form) || [];
  const currentIsAdvancedSelector = Form.useWatch('isAdvancedSelector', form) || false;
  const currentSelector = Form.useWatch('selector', form) || '';
  const query = currentIsAdvancedSelector
    ? currentSelector
    : SelectorService.getSelectorString(currentSelectorList, currentPseudoSelector);

  const onFieldsChange = useCallback(
    (changedFields: FieldData[]) => {
      const [field] = changedFields;

      const [fieldName = '', entry = 0, keyName = ''] = field.name as Array<string | number>;

      if (fieldName === 'isAdvancedSelector' && field.value) {
        form.setFieldsValue({
          selector: query,
        });
      }

      if (fieldName === 'selectorList') onSelectorChange();
      if (fieldName === 'assertionList') onChecksChange();

      if (fieldName === 'assertionList' && keyName === 'attribute' && field.value) {
        const list: TAssertion[] = form.getFieldValue('assertionList') || [];

        form.setFieldsValue({
          assertionList: list.map((assertionEntry, index) => {
            if (index === entry) {
              const {value = ''} = attributeList?.find((el: any) => el.key === list[index].attribute) || {};
              const isValid = typeof value === 'number' || !isEmpty(value);

              return {...assertionEntry, expected: isValid ? String(value) : ''};
            }

            return assertionEntry;
          }),
        });
      }
    },
    [attributeList, form, query]
  );

  return onFieldsChange;
};

export default useOnFieldsChange;
