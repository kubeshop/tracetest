import {useCallback, useState} from 'react';
import {useLazyGetJUnitByRunIdQuery} from 'redux/apis/Tracetest';

const useJUnitResult = () => {
  const [getJUnit] = useLazyGetJUnitByRunIdQuery();
  const [jUnit, setJUnit] = useState<string>('');

  const loadJUnit = useCallback(
    async (testId: string, runId: string) => {
      const data = await getJUnit({runId, testId}).unwrap();

      setJUnit(data);
    },
    [getJUnit]
  );

  return {jUnit, loadJUnit};
};

export default useJUnitResult;
