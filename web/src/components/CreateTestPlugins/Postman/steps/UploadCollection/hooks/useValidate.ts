import {useCallback} from 'react';
import Validator from '../../../../../../utils/Validator';
import {IRequestDetailsValues} from '../UploadCollection';

const useValidate = (
  onValidation: (isValid: boolean) => void,
  setTransientUrl: (value: ((prevState: string) => string) | string) => void
): ((changedValues: any, allValues: IRequestDetailsValues) => void) => {
  return useCallback(
    (changedValues, allValues: IRequestDetailsValues) => {
      onValidation(Validator.required(allValues.collectionTest) && Validator.required(allValues.collectionTest));
      if (changedValues.url) {
        setTransientUrl(changedValues.url);
      }
    },
    [onValidation]
  );
};

export default useValidate;
