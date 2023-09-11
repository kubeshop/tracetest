import {useCallback, useState} from 'react';
import TracetestAPI from 'redux/apis/Tracetest';

const {useLazyGetJUnitByRunIdQuery} = TracetestAPI.instance;

const useJUnitResult = () => {
  const [getJUnit] = useLazyGetJUnitByRunIdQuery();
  const [jUnit, setJUnit] = useState<string>('');

  const loadJUnit = useCallback(
    async (testId: string, runId: number) => {
      const data = await getJUnit({runId, testId}).unwrap();

      setJUnit(data);
    },
    [getJUnit]
  );

  return {jUnit, loadJUnit};
};

export default useJUnitResult;
