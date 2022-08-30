import {LoadingOutlined} from '@ant-design/icons';
import {Spin} from 'antd';

const antIcon = <LoadingOutlined style={{fontSize: 40}} spin />;

const LoadingSpinner = () => <Spin data-cy="loading-spinner" indicator={antIcon} />;

export default LoadingSpinner;
