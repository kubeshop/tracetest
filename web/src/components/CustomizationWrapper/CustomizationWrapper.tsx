import {useMemo} from 'react';
import CustomizationProvider from 'providers/Customization';

interface IProps {
  children: React.ReactNode;
}

const getComponent = <T,>(id: string, fallback: React.ComponentType<T>) => fallback;
const getIsAllowed = () => true;

const CustomizationWrapper = ({children}: IProps) => {
  const customizationProviderValue = useMemo(() => ({getComponent, getIsAllowed}), []);

  return <CustomizationProvider value={customizationProviderValue}>{children}</CustomizationProvider>;
};

export default CustomizationWrapper;
