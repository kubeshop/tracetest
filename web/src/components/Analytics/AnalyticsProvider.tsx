import {ReactElement, useEffect, useState, createContext, useCallback, useMemo} from 'react';
import {isEnabled, instance} from '../../entities/Analytics/Analytics.service';

export const Context = createContext({
  isEnabled,
  instance,
});

type TAnalyticsProviderProps = {
  children: ReactElement;
};

const AnalyticsProvider: React.FC<TAnalyticsProviderProps> = ({children}) => {
  const [isInitialized, setInitialized] = useState(!isEnabled);

  const initialize = useCallback(() => {
    instance.initialize().finally(() => setInitialized(true));
  }, []);

  useEffect(() => {
    if (isEnabled) initialize();
  }, [initialize]);

  const value = useMemo(() => ({isEnabled, instance}), []);

  return isInitialized ? <Context.Provider value={value}>{children}</Context.Provider> : <div />;
};

export default AnalyticsProvider;
