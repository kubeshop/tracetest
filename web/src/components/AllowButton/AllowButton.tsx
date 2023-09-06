import {Button, ButtonProps, Tooltip} from 'antd';
import {Operation, useCustomization} from 'providers/Customization/Customization.provider';

interface IProps extends ButtonProps {
  operation: Operation;
}

const AllowButton = ({operation, ...props}: IProps) => {
  const {getIsAllowed} = useCustomization();
  const isAllowed = getIsAllowed(operation);

  // the tooltip unmounts and remounts the children, detaching it from the DOM
  return isAllowed ? (
    <Button {...props} disabled={props.disabled} />
  ) : (
    <Tooltip title="You are not allowed to perform this operation">
      <Button {...props} disabled />
    </Tooltip>
  );
};

export default AllowButton;
