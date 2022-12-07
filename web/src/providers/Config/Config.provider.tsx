import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';

import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setUserPreference} from 'redux/slices/User.slice';
import {selectShouldDisplayConfigSetup} from 'redux/config/selectors';

interface IContext {
  shouldDisplayConfigSetup: boolean;
  skipConfigSetup(): void;
}

const Context = createContext<IContext>({
  shouldDisplayConfigSetup: true,
  skipConfigSetup: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useConfig = () => useContext(Context);

const ConfigProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const shouldDisplayConfigSetup = useAppSelector(selectShouldDisplayConfigSetup);

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
      shouldDisplayConfigSetup,
      skipConfigSetup,
    }),
    [shouldDisplayConfigSetup, skipConfigSetup]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default ConfigProvider;
