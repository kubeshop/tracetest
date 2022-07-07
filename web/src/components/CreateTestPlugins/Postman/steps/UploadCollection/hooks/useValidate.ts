import {useCallback} from 'react';
import {IRequestDetailsValues} from '../UploadCollection';

const useValidate = (
  onValidation: (isValid: boolean) => void
): ((changedValues: any, allValues: IRequestDetailsValues) => void) => {
  return useCallback(
    (changedValues: any, allValues: IRequestDetailsValues) => {
      // const isValid = Validator.required(allValues.url) && Validator.url(allValues.url);
      if (false) {
        console.log(allValues);
        console.log(changedValues);
      }
      onValidation(true);
    },
    [onValidation]
  );
};

export default useValidate;
