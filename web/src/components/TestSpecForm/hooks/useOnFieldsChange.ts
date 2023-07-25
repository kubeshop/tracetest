import {useCallback} from 'react';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';

interface FieldData {
  name: string | number | (string | number)[];
  value?: any;
  touched?: boolean;
  validating?: boolean;
  errors?: string[];
}

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
