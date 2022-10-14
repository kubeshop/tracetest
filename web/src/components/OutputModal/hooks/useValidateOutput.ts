import {useCallback, useState} from 'react';
import {TOutput} from 'types/Output.types';
import {SupportedEditors} from 'constants/Editor.constants';
import useEditorValidate from 'components/Editor/hooks/useEditorValidate';

interface IProps {
  spanIdList: string[];
}

const useValidateOutput = ({spanIdList}: IProps) => {
  const [isFormValid, setIsFormValid] = useState(false);

  const getIsValidSelector = useEditorValidate();

  const onValidate = useCallback(
    async (changedValues: any, {source, selector, attribute}: TOutput) => {
      const isBaseValid = Boolean(source && attribute);

      if (source === 'trace') {
        setIsFormValid(
          isBaseValid &&
            spanIdList.length === 1 &&
            Boolean(selector) &&
            getIsValidSelector(SupportedEditors.Selector, selector || '')
        );
      } else {
        setIsFormValid(isBaseValid);
      }
    },
    [getIsValidSelector, spanIdList.length]
  );

  return {isValid: isFormValid, onValidate};
};

export default useValidateOutput;
