import {CloseCircleFilled} from '@ant-design/icons';
import {Typography} from 'antd';
import * as S from './ErrorBoundary.styled';

const ErrorBoundary: React.FC = () => {
  return (
    <S.Container>
      <CloseCircleFilled style={{color: 'red', fontSize: 32}} />
      <Typography.Title level={2}>Something went wrong!</Typography.Title>
    </S.Container>
  );
};

export default ErrorBoundary;
