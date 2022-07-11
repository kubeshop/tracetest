import {useCallback, useState} from 'react';
import {TriggerTypes} from 'constants/Test.constants';
import TestService from 'services/Test.service';
import {TDraftTest} from 'types/Test.types';

interface IProps {
  isDefaultValid?: boolean;
  type: TriggerTypes;
  isBasicDetails?: boolean;
}

const useValidateTestDraft = ({isDefaultValid = false, type, isBasicDetails = false}: IProps) => {
  const [isFormValid, setIsFormValid] = useState(isDefaultValid);

  const onValidate = useCallback(
    async (changedValues: any, draft: TDraftTest) => {
      const isValid = await TestService.validateDraft(type, draft, isBasicDetails);

      setIsFormValid(isValid);
    },
    [isBasicDetails, type]
  );

  return {isValid: isFormValid, setIsValid: setIsFormValid, onValidate};
};

export default useValidateTestDraft;
