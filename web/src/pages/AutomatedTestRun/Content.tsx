import {useEffect} from 'react';
import {useSearchParams} from 'react-router-dom';
import useAutomatedTestRun from './hooks/useAutomatedTestRun';
import TestContent from '../Test/Content';

const Content = () => {
  const [query] = useSearchParams();
  const onAutomatedRun = useAutomatedTestRun(query);

  useEffect(() => {
    onAutomatedRun();
  }, [onAutomatedRun]);

  return <TestContent />;
};

export default Content;
