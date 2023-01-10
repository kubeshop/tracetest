import {useCallback, useState} from 'react';
import {TDraftVariables} from 'types/Variables.types';

const useValidateVariablesDraft = () => {
  const [isFormValid, setIsFormValid] = useState(false);

  const onValidate = useCallback((changedValues: any, {variables}: TDraftVariables) => {
    const found = Object.values(variables).find(value => !value);
    const hasAMissingVariable = found !== undefined;

    setIsFormValid(!hasAMissingVariable);
  }, []);

  return {isValid: isFormValid, onValidate};
};

export default useValidateVariablesDraft;
