import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import {useGetDataStoreConfigQuery} from 'redux/apis/TraceTest.api';
import {ConfigMode, TDataStoreConfig} from 'types/Config.types';
import UserSelectors from 'selectors/User.selectors';
import DataStoreConfig from 'models/DataStoreConfig.model';

interface IContext {
  dataStoreConfig: TDataStoreConfig;
  isLoading: boolean;
  isError: boolean;
  skipConfigSetup(): void;
  shouldDisplayConfigSetup: boolean;
}

const Context = createContext<IContext>({
  dataStoreConfig: DataStoreConfig({}),
  skipConfigSetup: noop,
  isLoading: false,
  isError: false,
  shouldDisplayConfigSetup: false,
});

interface IProps {
  children: React.ReactNode;
}

export const useDataStoreConfig = () => useContext(Context);

const DataStoreConfigProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {data: dataStoreConfig = DataStoreConfig({}), isLoading, isError} = useGetDataStoreConfigQuery({});
  const initConfigSetup = useAppSelector(state => UserSelectors.selectUserPreference(state, 'initConfigSetup'));

  const shouldDisplayConfigSetup = Boolean(initConfigSetup) && dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

  const skipConfigSetup = useCallback(() => {
    dispatch(
      setUserPreference({
        key: 'initConfigSetup',
        value: false,
      })
    );
  }, [dispatch]);

  const value = useMemo<IContext>(
    () => ({
      dataStoreConfig,
      isLoading,
      isError,
      skipConfigSetup,
      shouldDisplayConfigSetup,
    }),
    [dataStoreConfig, isError, isLoading, shouldDisplayConfigSetup, skipConfigSetup]
  );

  return isLoading ? null : <Context.Provider value={value}>{children}</Context.Provider>;
};

export default DataStoreConfigProvider;
