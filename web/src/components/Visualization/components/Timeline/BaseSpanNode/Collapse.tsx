import * as S from '../Timeline.styled';

interface IProps {
  id: string;
  isCollapsed?: boolean;
  onCollapse(id: string): void;
  totalChildren: number;
}

const collapsedShape = 'M18.629 15.997l-7.083-7.081L13.462 7l8.997 8.997L13.457 25l-1.916-1.916z';
const expandedShape = 'M16.003 18.626l7.081-7.081L25 13.46l-8.997 8.998-9.003-9 1.917-1.916z';

const Collapse = ({id, isCollapsed, onCollapse, totalChildren}: IProps) => (
  <S.GroupCollapse
    left={-5}
    onClick={event => {
      event.stopPropagation();
      onCollapse(id);
    }}
    top={0}
  >
    <S.CircleArrow cx={10} cy={10} r={8} />
    <S.PathArrow d={isCollapsed ? collapsedShape : expandedShape} />
    <S.CircleNumber cx={26} cy={10} r={7} />
    <S.TextNumber dominantBaseline="hanging" textAnchor="middle" x={26} y={5}>
      {totalChildren}
    </S.TextNumber>
  </S.GroupCollapse>
);

export default Collapse;
