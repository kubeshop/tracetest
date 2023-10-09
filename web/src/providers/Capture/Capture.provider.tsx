import noop from 'lodash/noop';
import {createContext, useContext} from 'react';

interface IContext {
  identify(): void;
  pageView(): void;
}

export const Context = createContext<IContext>({
  identify: noop,
  pageView: noop,
});

export const useCapture = () => useContext(Context);

interface IProps {
  children: React.ReactNode;
  value: IContext;
}

const CaptureProvider = ({children, value}: IProps) => {
  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CaptureProvider;
