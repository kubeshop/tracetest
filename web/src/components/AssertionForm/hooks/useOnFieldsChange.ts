import {FormInstance} from 'antd';
import {FieldData} from 'antd/node_modules/rc-field-form/es/interface';
import {isEmpty} from 'lodash';
import {useCallback} from 'react';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {TAssertion} from 'types/Assertion.types';
import {TSpanFlatAttribute} from '../../../types/Span.types';
import {IValues} from '../AssertionForm';

const {onChecksChange, onSelectorChange} = CreateAssertionModalAnalyticsService;

interface IProps {
  form: FormInstance<IValues>;
  attributeList: TSpanFlatAttribute[];
}

const useOnFieldsChange = ({form, attributeList}: IProps) => {
  return useCallback(
    (changedFields: FieldData[]) => {
      const [field] = changedFields;

      if (field?.name) {
        const [fieldName = '', entry = 0, keyName = ''] = field.name as Array<string | number>;

        if (fieldName === 'selector') onSelectorChange();
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
      }
    },
    [attributeList, form]
  );
};

export default useOnFieldsChange;
