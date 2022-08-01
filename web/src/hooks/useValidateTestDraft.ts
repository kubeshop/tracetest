import {useCallback, useState} from 'react';
import TestService from 'services/Test.service';
import {TDraftTest} from 'types/Test.types';
import {SupportedPlugins} from 'constants/Plugins.constants';

interface IProps {
  isDefaultValid?: boolean;
  pluginName: SupportedPlugins;
  isBasicDetails?: boolean;
}

const useValidateTestDraft = ({isDefaultValid = false, pluginName, isBasicDetails = false}: IProps) => {
  const [isFormValid, setIsFormValid] = useState(isDefaultValid);

  const onValidate = useCallback(
    async (changedValues: any, draft: TDraftTest) => {
      const isValid = await TestService.validateDraft(pluginName, draft, isBasicDetails);

      setIsFormValid(isValid);
    },
    [pluginName, isBasicDetails]
  );

  return {isValid: isFormValid, setIsValid: setIsFormValid, onValidate};
};

export default useValidateTestDraft;
