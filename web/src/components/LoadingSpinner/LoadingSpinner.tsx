import {Spin} from 'antd';
import {LoadingOutlined} from '@ant-design/icons';

const antIcon = <LoadingOutlined style={{fontSize: 40}} spin />;

const LoadingSpinner: React.FC = () => <Spin indicator={antIcon} />;

export default LoadingSpinner;
