import {useCallback, useEffect, useState} from 'react';
import {useSearchParams} from 'react-router-dom';

const useQueryTabs = (defaultTab = '') => {
  const [query, setQuery] = useSearchParams();
  const [activeKey, setActiveKey] = useState(query.get('tab') || defaultTab);

  useEffect(() => {
    const tab = query.get('tab');
    if (tab) setActiveKey(tab);
  }, [query]);

  const handleChange = useCallback(
    (newTab: string) => {
      setQuery([['tab', newTab]]);
    },
    [setQuery]
  );

  return [activeKey, handleChange] as const;
};

export default useQueryTabs;
