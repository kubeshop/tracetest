import {useCustomization} from './Customization.provider';

const withCustomization = <P extends object>(Component: React.ComponentType<P>, id: string) => {
  const WrappedComponent = (props: P) => {
    const {getComponent} = useCustomization();
    const CustomizedComponent = getComponent(id, Component);

    return <CustomizedComponent {...props} />;
  };

  return WrappedComponent;
};

export default withCustomization;
