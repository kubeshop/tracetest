import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import {useGetDataStoresQuery} from 'redux/apis/TraceTest.api';
import {ConfigMode} from 'types/DataStore.types';
import UserSelectors from 'selectors/User.selectors';
import DataStoreConfig from 'models/DataStoreConfig.model';

interface IContext {
  dataStoreConfig: DataStoreConfig;
  isLoading: boolean;
  isFetching: boolean;
  isError: boolean;
  skipConfigSetup(): void;
  skipConfigSetupFromTest(): void;
  shouldDisplayConfigSetup: boolean;
  shouldDisplayConfigSetupFromTest: boolean;
}

const Context = createContext<IContext>({
  dataStoreConfig: DataStoreConfig([]),
  skipConfigSetup: noop,
  skipConfigSetupFromTest: noop,
  isLoading: false,
  isFetching: false,
  isError: false,
  shouldDisplayConfigSetup: false,
  shouldDisplayConfigSetupFromTest: false,
});

interface IProps {
  children: React.ReactNode;
}

export const useDataStoreConfig = () => useContext(Context);

const DataStoreConfigProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {data: dataStoreConfig = DataStoreConfig([]), isLoading, isError, isFetching} = useGetDataStoresQuery({});
  const initConfigSetup = useAppSelector(state => UserSelectors.selectUserPreference(state, 'initConfigSetup'));
  const initConfigSetupFromTest = useAppSelector(state =>
    UserSelectors.selectUserPreference(state, 'initConfigSetupFromTest')
  );

  const shouldDisplayConfigSetup = !!initConfigSetup && dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;
  const shouldDisplayConfigSetupFromTest =
    !!initConfigSetupFromTest && dataStoreConfig.mode === ConfigMode.NO_TRACING_MODE;

  const skipConfigSetup = useCallback(() => {
    dispatch(
      setUserPreference({
        key: 'initConfigSetup',
        value: false,
      })
    );
  }, [dispatch]);

  const skipConfigSetupFromTest = useCallback(() => {
    dispatch(
      setUserPreference({
        key: 'initConfigSetupFromTest',
        value: false,
      })
    );
  }, [dispatch]);

  const value = useMemo<IContext>(
    () => ({
      dataStoreConfig,
      isLoading,
      isFetching,
      isError,
      skipConfigSetup,
      skipConfigSetupFromTest,
      shouldDisplayConfigSetup,
      shouldDisplayConfigSetupFromTest,
    }),
    [
      dataStoreConfig,
      isError,
      isLoading,
      isFetching,
      shouldDisplayConfigSetup,
      shouldDisplayConfigSetupFromTest,
      skipConfigSetup,
      skipConfigSetupFromTest,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default DataStoreConfigProvider;
