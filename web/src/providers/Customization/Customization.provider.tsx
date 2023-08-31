import {createContext, useContext} from 'react';

export enum Operation {
  Configure = 'configure',
  Edit = 'edit',
  View = 'view',
}

interface IContext {
  getComponent<T>(name: string, fallback: React.ComponentType<T>): React.ComponentType<T>;
  getIsAllowed(operation: Operation): boolean;
}

export const Context = createContext<IContext>({
  getComponent: (name, fallback) => fallback,
  getIsAllowed: () => true,
});

export const useCustomization = () => useContext(Context);

interface IProps {
  children: React.ReactNode;
  value: IContext;
}

const CustomizationProvider = ({children, value}: IProps) => {
  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default CustomizationProvider;
