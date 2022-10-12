import {useCallback, useState} from 'react';
import {TOutput} from 'types/Output.types';
import SelectorService from 'services/Selector.service';

interface IProps {
  spanIdList: string[];
}

const useValidateOutput = ({spanIdList}: IProps) => {
  const [isFormValid, setIsFormValid] = useState(false);

  const onValidate = useCallback(
    async (changedValues: any, {source, selector, attribute}: TOutput) => {
      const isBaseValid = Boolean(source && attribute);

      if (source === 'trace') {
        setIsFormValid(
          isBaseValid &&
            spanIdList.length === 1 &&
            Boolean(selector) &&
            SelectorService.getIsValidSelector(selector || '')
        );
      } else {
        setIsFormValid(isBaseValid);
      }
    },
    [spanIdList.length]
  );

  return {isValid: isFormValid, onValidate};
};

export default useValidateOutput;
