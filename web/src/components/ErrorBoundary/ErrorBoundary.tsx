import {Typography} from 'antd';
import * as S from './ErrorBoundary.styled';

interface IErrorBoundaryProps {
  error: Error;
}

const ErrorBoundary: React.FC<IErrorBoundaryProps> = ({error}) => {
  return (
    <S.Container>
      <S.Icon />
      <Typography.Title level={1}>Something went wrong!</Typography.Title>
      <S.Content>{error.toString()}</S.Content>
    </S.Container>
  );
};

export default ErrorBoundary;
