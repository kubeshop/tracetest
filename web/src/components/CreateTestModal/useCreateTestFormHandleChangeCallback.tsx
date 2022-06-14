import {useCallback} from 'react';
import Validator from '../../utils/Validator';
import {ICreateTestValues} from './CreateTestForm';

export const useCreateTestFormHandleChangeCallback = (
  onValidation: (isValid: boolean) => void
): ((changedValues: any, allValues: ICreateTestValues) => void) => {
  const authValidation = (testValues: ICreateTestValues): boolean => {
    switch (testValues.auth?.type) {
      case 'apiKey':
        return Validator.required(testValues.auth?.apiKey?.key) && Validator.required(testValues.auth?.apiKey?.value);
      case 'basic':
        return (
          Validator.required(testValues.auth?.basic?.username) && Validator.required(testValues.auth?.basic?.password)
        );
      case 'bearer':
        return Validator.required(testValues.auth?.bearer?.token);
      default:
        return true;
    }
  };
  return useCallback(
    (changedValues: any, allValues: ICreateTestValues) => {
      const isValid =
        Validator.required(allValues.name) &&
        Validator.required(allValues.url) &&
        Validator.url(allValues.url) &&
        authValidation(allValues);
      onValidation(isValid);
    },
    [onValidation]
  );
};
