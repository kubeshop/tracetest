import * as S from './Collapse.styled';
import CollapseIcon from './CollapseIcon';

const Collapse: React.FC = ({children}) => {
  return (
    <S.StyledCollapse expandIcon={({isActive = false}) => <CollapseIcon isCollapsed={isActive} />}>
      {children}
    </S.StyledCollapse>
  );
};

export default Collapse;
