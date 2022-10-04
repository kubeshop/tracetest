import {useCallback} from 'react';
import TestService from 'services/Test.service';
import {TDraftTest} from 'types/Test.types';
import {SupportedPlugins} from 'constants/Common.constants';

interface IProps {
  setIsValid(isValid: boolean): void;
  pluginName: SupportedPlugins;
  isBasicDetails?: boolean;
}

const useValidateTestDraft = ({pluginName, isBasicDetails = false, setIsValid}: IProps) => {
  const onValidate = useCallback(
    async (changedValues: any, draft: TDraftTest) => {
      const isValid = await TestService.validateDraft(pluginName, draft, isBasicDetails);

      setIsValid(isValid);
    },
    [pluginName, isBasicDetails, setIsValid]
  );

  return onValidate;
};

export default useValidateTestDraft;
