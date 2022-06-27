import {useCallback} from 'react';
import Validator from 'utils/Validator';
import {IRequestDetailsValues} from '../RequestDetails';

const useValidate = (
  onValidation: (isValid: boolean) => void
): ((changedValues: any, allValues: IRequestDetailsValues) => void) => {
  const authValidation = (testValues: IRequestDetailsValues): boolean => {
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
    (changedValues: any, allValues: IRequestDetailsValues) => {
      const isValid = Validator.required(allValues.url) && Validator.url(allValues.url) && authValidation(allValues);
      onValidation(isValid);
    },
    [onValidation]
  );
};

export default useValidate;
