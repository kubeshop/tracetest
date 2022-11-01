import {useCallback, useState} from 'react';

import {SupportedEditors} from 'constants/Editor.constants';
import EditorService from 'services/Editor.service';
import {TTestOutput} from 'types/TestOutput.types';

interface IProps {
  spanIdList: string[];
}

const useValidateOutput = ({spanIdList}: IProps) => {
  const [isFormValid, setIsFormValid] = useState(false);

  const onValidate = useCallback(
    (changedValues: any, {name, selector, value}: TTestOutput) => {
      setIsFormValid(
        Boolean(name) &&
          Boolean(selector) &&
          Boolean(value) &&
          spanIdList.length === 1 &&
          EditorService.getIsQueryValid(SupportedEditors.Selector, selector || '')
      );
    },
    [spanIdList.length]
  );

  return {isValid: isFormValid, onValidate};
};

export default useValidateOutput;
