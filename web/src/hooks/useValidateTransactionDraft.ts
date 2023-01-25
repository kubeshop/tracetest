import {useCallback} from 'react';
import {TDraftTransaction} from 'types/Transaction.types';

interface IProps {
  setIsValid(isValid: boolean): void;
}

const useValidateTransactionDraft = ({setIsValid}: IProps) => {
  const onValidate = useCallback(
    async (changedValues: any, draft: TDraftTransaction) => {
      const isValid = !!draft.name;

      setIsValid(isValid);
    },
    [setIsValid]
  );

  return onValidate;
};

export default useValidateTransactionDraft;
