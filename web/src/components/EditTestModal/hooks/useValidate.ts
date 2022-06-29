import {useCallback} from 'react';
import Validator from 'utils/Validator';
import {TEditTest} from '../EditTestModal';

const useValidate = (
  onValidation: (isValid: boolean) => void
): ((changedValues: any, allValues: TEditTest) => void) => {
  const authValidation = ({auth}: TEditTest): boolean => {
    switch (auth?.type) {
      case 'apiKey':
        return Validator.required(auth?.apiKey?.key) && Validator.required(auth?.apiKey?.value);
      case 'basic':
        return Validator.required(auth?.basic?.username) && Validator.required(auth?.basic?.password);
      case 'bearer':
        return Validator.required(auth?.bearer?.token);
      default:
        return true;
    }
  };
  return useCallback(
    (changedValues: any, allValues: TEditTest) => {
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

export default useValidate;
