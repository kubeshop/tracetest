import {Dispatch, SetStateAction, useCallback} from 'react';
import Validator from 'utils/Validator';
import {IUploadCollectionValues} from '../UploadCollection';

const useValidate = (
  onValidation: (isValid: boolean) => void,
  setTransientUrl: Dispatch<SetStateAction<string>>
): ((changedValues: any, allValues: IUploadCollectionValues) => void) => {
  return useCallback(
    (changedValues, allValues: IUploadCollectionValues) => {
      onValidation(Validator.required(allValues.collectionTest) && Validator.required(allValues.collectionTest));
      if (changedValues.url) {
        setTransientUrl(changedValues.url);
      }
    },
    [onValidation, setTransientUrl]
  );
};

export default useValidate;
