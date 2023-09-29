import {Skeleton} from 'antd';
import * as S from './ResourceCard.styled';

interface IProps {
  children: React.ReactNode;
  hasMoreRuns: boolean;
  hasRuns: boolean;
  isCollapsed: boolean;
  isLoading: boolean;
  onViewAll(): void;
  resourcePath: string;
}

const ResourceCardRuns = ({
  children,
  hasMoreRuns,
  hasRuns,
  isCollapsed,
  isLoading,
  onViewAll,
  resourcePath,
}: IProps) => {
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
          <S.CustomLink data-cy="test-details-link" onClick={onViewAll} to={resourcePath}>
            View all runs
          </S.CustomLink>
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
