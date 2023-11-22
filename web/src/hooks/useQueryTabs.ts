import {useCallback, useEffect, useState} from 'react';
import {useSearchParams} from 'react-router-dom';

const useQueryTabs = (defaultTab = '', paramName = 'tab') => {
  const [query, setQuery] = useSearchParams();
  const [activeKey, setActiveKey] = useState(query.get(paramName) || defaultTab);

  useEffect(() => {
    const tab = query.get(paramName);
    if (tab) setActiveKey(tab);
  }, [paramName, query]);

  const handleChange = useCallback(
    (newTab: string) => {
      setQuery([[paramName, newTab]]);
    },
    [paramName, setQuery]
  );

  return [activeKey, handleChange] as const;
};

export default useQueryTabs;
