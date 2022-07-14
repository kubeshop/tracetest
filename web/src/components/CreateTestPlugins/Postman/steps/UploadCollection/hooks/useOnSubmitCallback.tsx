import {useCallback} from 'react';
import {useCreateTest} from '../../../../../../providers/CreateTest/CreateTest.provider';
import {IUploadCollectionValues} from '../UploadCollection';

export function useOnSubmitCallback() {
  const {onNext} = useCreateTest();
  return useCallback(
    ({collectionFile, envFile, collectionTest, requests, variables, ...values}: IUploadCollectionValues) => {
      // eslint-disable-next-line no-console
      console.log(collectionFile, envFile, collectionTest, requests, variables);
      onNext({serviceUnderTest: {triggerSettings: {http: values}}});
    },
    [onNext]
  );
}
