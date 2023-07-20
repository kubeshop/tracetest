import {FieldData} from 'rc-field-form/es/interface';
import {useCallback} from 'react';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';

const {onChecksChange, onSelectorChange} = CreateAssertionModalAnalyticsService;

const useOnFieldsChange = () => {
  return useCallback((changedFields: FieldData[]) => {
    const [field] = changedFields;

    if (field?.name) {
      const [fieldName = ''] = field.name as Array<string | number>;

      if (fieldName === 'selector') onSelectorChange();
      if (fieldName === 'assertions') onChecksChange();
    }
  }, []);
};

export default useOnFieldsChange;
