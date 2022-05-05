import {CloseCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import * as S from './ErrorBoundary.styled';

interface IErrorBoundaryProps {
  error: Error;
}

const ErrorBoundary: React.FC<IErrorBoundaryProps> = ({error}) => {
  return (
    <S.Container>
      <CloseCircleFilled style={{color: 'red', fontSize: 32}} />
      <Typography.Title level={2}>Something went wrong!</Typography.Title>
      <div style={{display: 'flex', maxWidth: '800px', padding: '24px'}}>{error.toString()}</div>
    </S.Container>
  );
};

export default ErrorBoundary;
