import {useCallback} from 'react';
import Validator from 'utils/Validator';
import {IBasicDetailsValues} from '../BasicDetails';

const useValidate = (
  onValidation: (isValid: boolean) => void
): ((changedValues: any, allValues: IBasicDetailsValues) => void) => {
  return useCallback(
    (changedValues: any, allValues: IBasicDetailsValues) => {
      const isValid = Validator.required(allValues.name) && Validator.required(allValues.description);
      onValidation(isValid);
    },
    [onValidation]
  );
};

export default useValidate;