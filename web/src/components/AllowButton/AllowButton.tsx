import {Button, ButtonProps, Tooltip, Typography} from 'antd';
import {capitalize} from 'lodash';
import {Operation, useCustomization} from 'providers/Customization/Customization.provider';
import * as S from './AllowButton.styled';

interface IProps extends ButtonProps {
  operation: Operation;
  ButtonComponent?: React.ComponentType<ButtonProps>;
}

const AllowButton = ({operation, ButtonComponent, ...props}: IProps) => {
  const {getIsAllowed, getRole} = useCustomization();
  const isAllowed = getIsAllowed(operation);
  const BtnComponent = ButtonComponent || Button;
  const role = getRole();

  // the tooltip unmounts and remounts the children, detaching it from the DOM
  return isAllowed ? (
    <BtnComponent {...props} disabled={props.disabled} />
  ) : (
    <Tooltip
      placement="topRight"
      overlay={
        <>
          <Typography.Title level={3}>
            <S.Warning color="yellow" /> Limited Access
          </Typography.Title>
          <Typography.Paragraph>
            Your current role group (<b>{capitalize(role)}</b>) has limited access to this environment. please contact
            the environment administrator for assistance.
          </Typography.Paragraph>
        </>
      }
      overlayStyle={{minWidth: '420px'}}
    >
      <Button {...props} disabled />
    </Tooltip>
  );
};

export default AllowButton;
