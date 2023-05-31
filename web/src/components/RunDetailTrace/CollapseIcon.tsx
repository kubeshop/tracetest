import * as S from './LintResults.styled';

interface IProps {
  isCollapsed: boolean;
}

const CollapseIcon = ({isCollapsed}: IProps) => {
  return (
    <S.CollapseIconContainer>{isCollapsed ? <S.DownCollapseIcon /> : <S.UpCollapseIcon />}</S.CollapseIconContainer>
  );
};

export default CollapseIcon;
