import {useCallback} from 'react';
import {TDraftTransaction} from 'types/Transaction.types';

interface IProps {
  setIsValid(isValid: boolean): void;
}

const useValidateTransactionDraft = ({setIsValid}: IProps) => {
  const onValidate = useCallback(
    async (changedValues: any, draft: TDraftTransaction) => {
      const isValid = Boolean(draft.name && draft.description);

      setIsValid(isValid);
    },
    [setIsValid]
  );

  return onValidate;
};

export default useValidateTransactionDraft;
