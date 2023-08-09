import {useCallback} from 'react';
import {TDraftTestSuite} from 'types/TestSuite.types';

interface IProps {
  setIsValid(isValid: boolean): void;
}

const useValidateTestSuiteDraft = ({setIsValid}: IProps) => {
  const onValidate = useCallback(
    async (changedValues: any, draft: TDraftTestSuite) => {
      const isValid = !!draft.name;

      setIsValid(isValid);
    },
    [setIsValid]
  );

  return onValidate;
};

export default useValidateTestSuiteDraft;
