import {Button, ButtonProps, Tooltip} from 'antd';
import {Operation, useCustomization} from 'providers/Customization/Customization.provider';

interface IProps extends ButtonProps {
  operation: Operation;
  ButtonComponent?: React.ComponentType<ButtonProps>;
}

const AllowButton = ({operation, ButtonComponent, ...props}: IProps) => {
  const {getIsAllowed} = useCustomization();
  const isAllowed = getIsAllowed(operation);
  const BtnComponent = ButtonComponent || Button;

  // the tooltip unmounts and remounts the children, detaching it from the DOM
  return isAllowed ? (
    <BtnComponent {...props} disabled={props.disabled} />
  ) : (
    <Tooltip title="You are not allowed to perform this operation">
      <Button {...props} disabled />
    </Tooltip>
  );
};

export default AllowButton;
