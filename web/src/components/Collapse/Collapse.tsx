import * as S from './Collapse.styled';
import CollapseIcon from './CollapseIcon';

interface IProps {
  children: React.ReactNode;
  onChange?(): void;
}

const Collapse = ({children, onChange}: IProps) => {
  return (
    <S.StyledCollapse expandIcon={({isActive = false}) => <CollapseIcon isCollapsed={isActive} />} onChange={onChange}>
      {children}
    </S.StyledCollapse>
  );
};

export default Collapse;
