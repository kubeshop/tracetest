import * as S from './Collapse.styled';

interface IProps {
  isCollapsed: boolean;
}

const CollapseIcon = ({isCollapsed}: IProps) => {
  return (
    <S.CollapseIconContainer>{isCollapsed ? <S.UpCollapseIcon /> : <S.DownCollapseIcon />}</S.CollapseIconContainer>
  );
};

export default CollapseIcon;
