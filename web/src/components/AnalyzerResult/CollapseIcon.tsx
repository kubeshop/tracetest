import * as S from './AnalyzerResult.styled';

interface IProps {
  isCollapsed: boolean;
}

const CollapseIcon = ({isCollapsed}: IProps) => {
  return (
    <S.CollapseIconContainer>{isCollapsed ? <S.DownCollapseIcon /> : <S.UpCollapseIcon />}</S.CollapseIconContainer>
  );
};

export default CollapseIcon;
