import {Skeleton} from 'antd';
import * as S from './ResourceCard.styled';

interface IProps {
  children: React.ReactNode;
  hasMoreRuns: boolean;
  hasRuns: boolean;
  isCollapsed: boolean;
  isLoading: boolean;
  onViewAll(): void;
}

const ResourceCardRuns = ({children, hasMoreRuns, hasRuns, isCollapsed, isLoading, onViewAll}: IProps) => {
  if (isCollapsed) return null;

  return (
    <S.RunsContainer>
      {isLoading && (
        <S.LoadingContainer direction="vertical">
          <Skeleton.Input active block size="small" />
          <Skeleton.Input active block size="small" />
          <Skeleton.Input active block size="small" />
        </S.LoadingContainer>
      )}

      {hasRuns && children}

      {hasMoreRuns && (
        <S.FooterContainer>
          <S.Link data-cy="test-details-link" onClick={onViewAll}>
            View all runs
          </S.Link>
        </S.FooterContainer>
      )}

      {!hasRuns && !isLoading && (
        <S.EmptyStateContainer>
          <S.EmptyStateIcon />
          <S.Text disabled>No Runs</S.Text>
        </S.EmptyStateContainer>
      )}
    </S.RunsContainer>
  );
};

export default ResourceCardRuns;
