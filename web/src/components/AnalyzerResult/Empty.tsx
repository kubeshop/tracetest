import {Button} from 'antd';
import {Link} from 'react-router-dom';
import * as S from './AnalyzerResult.styled';

const Empty = () => {
  return (
    <S.EmptyContainer data-cy="empty-analyzer-results">
      <S.EmptyIcon />
      <S.EmptyTitle>There are no Analyzer results yet</S.EmptyTitle>
      <S.EmptyText>Please configure the Tracetest Analyzer settings to see result</S.EmptyText>
      <S.ConfigureButtonContainer>
        <Link to="/settings?tab=analyzer">
          <Button ghost type="primary">
            Configure
          </Button>
        </Link>
      </S.ConfigureButtonContainer>
    </S.EmptyContainer>
  );
};

export default Empty;
