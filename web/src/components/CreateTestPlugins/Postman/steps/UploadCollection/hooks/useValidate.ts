import {Dispatch, SetStateAction, useCallback} from 'react';
import Validator from 'utils/Validator';
import {IRequestDetailsValues} from '../UploadCollection';

const useValidate = (
  onValidation: (isValid: boolean) => void,
  setTransientUrl: Dispatch<SetStateAction<string>>
): ((changedValues: any, allValues: IRequestDetailsValues) => void) => {
  return useCallback(
    (changedValues, allValues: IRequestDetailsValues) => {
      onValidation(Validator.required(allValues.collectionTest) && Validator.required(allValues.collectionTest));
      if (changedValues.url) {
        setTransientUrl(changedValues.url);
      }
    },
    [onValidation, setTransientUrl]
  );
};

export default useValidate;
