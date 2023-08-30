import {useMemo} from 'react';
import CustomizationProvider from 'providers/CustomizationProvider/Customization.provider';

interface IProps {
  children: React.ReactNode;
}

const CustomizationWrapper = ({children}: IProps) => {
  const getComponent = <T,>(id: string, fallback: React.ComponentType<T>) => fallback;
  const getIsAllowed = () => true;
  const customizationProviderValue = useMemo(() => ({getComponent, getIsAllowed}), []);

  return <CustomizationProvider value={customizationProviderValue}>{children}</CustomizationProvider>;
};

export default CustomizationWrapper;
