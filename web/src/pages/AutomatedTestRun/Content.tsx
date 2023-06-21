import {useEffect, useState} from 'react';
import {useSearchParams} from 'react-router-dom';
import useAutomatedTestRun from './hooks/useAutomatedTestRun';
import TestContent from '../Test/Content';

const Content = () => {
  const [query] = useSearchParams();
  const [hasTriggered, setHasTriggered] = useState(false);
  const onAutomatedRun = useAutomatedTestRun(query);

  useEffect(() => {
    if (!hasTriggered) {
      setHasTriggered(true);
      onAutomatedRun();
    }
  }, [hasTriggered, onAutomatedRun]);

  return <TestContent />;
};

export default Content;
