import {Button, ButtonProps, Tooltip} from 'antd';
import {Operation, useCustomization} from 'providers/Customization/Customization.provider';

interface IProps extends ButtonProps {
  operation: Operation;
}

const AllowButton = ({operation, ...props}: IProps) => {
  const {getIsAllowed} = useCustomization();
  const isAllowed = getIsAllowed(operation);

  return (
    <Tooltip title={!isAllowed ? 'You are not allowed to perform this operation' : ''}>
      <Button {...props} disabled={!isAllowed} />
    </Tooltip>
  );
};

export default AllowButton;
