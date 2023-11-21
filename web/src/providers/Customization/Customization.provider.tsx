import {createContext, useContext} from 'react';

export enum Operation {
  Configure = 'configure',
  Edit = 'edit',
  View = 'view',
}

export enum Flag {
  IsAgentDataStoreEnabled = 'isAgentDataStoreEnabled',
  IsLocalModeEnabled = 'isLocalModeEnabled',
}

interface IContext {
  getComponent<T>(name: string, fallback: React.ComponentType<T>): React.ComponentType<T>;
  getFlag(flag: Flag): boolean;
  getIsAllowed(operation: Operation): boolean;
  getRole(): string;
}

export const Context = createContext<IContext>({
  getComponent: (name, fallback) => fallback,
  getFlag: () => true,
  getIsAllowed: () => true,
  getRole: () => '',
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
