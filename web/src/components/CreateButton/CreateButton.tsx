import {Button, ButtonProps} from 'antd';
import {withCustomization} from 'providers/Customization';

interface IProps extends ButtonProps {}

const CreateButton = ({...props}: IProps) => <Button {...props} />;

export default withCustomization(CreateButton, 'createButton');
