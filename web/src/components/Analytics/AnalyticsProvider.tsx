import {ReactElement, useEffect, useState, createContext, useCallback, useMemo} from 'react';
import GA4React from 'ga-4-react';

const {analyticsEnabled = 'true', measurementId = 'G-ZP277L2M37'} = window.ENV || {};

export const instance = new GA4React(measurementId);

const isEnabled = analyticsEnabled === 'true';

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
