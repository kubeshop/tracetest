import {useCallback} from 'react';
import {IValues} from '../TestSpecForm';

interface IProps {
  setIsValid(isValid: boolean): void;
}

const useOnValuesChange = ({setIsValid}: IProps) => {
  return useCallback(
    (_: any, {assertions = []}: IValues) => {
      const isValid = !assertions.find(assertion => !assertion?.left || !assertion?.right);
      setIsValid(!!assertions.length && isValid);
    },
    [setIsValid]
  );
};

export default useOnValuesChange;
