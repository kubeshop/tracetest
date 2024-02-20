import Ticks from './Ticks/Ticks';
import * as S from './Timeline.styled';
import VerticalResizer from './VerticalResizer/VerticalResizer';

const NUM_TICKS = 5;

interface IProps {
  duration: number;
  nameColumnWidth: number;
  onNameColumnWidthChange(width: number): void;
}

const Header = ({duration, nameColumnWidth, onNameColumnWidthChange}: IProps) => (
  <S.HeaderRow style={{gridTemplateColumns: `${nameColumnWidth * 100}% 1fr`}}>
    <S.Col>
      <S.HeaderContent>
        <S.HeaderTitle level={3}>Span</S.HeaderTitle>
      </S.HeaderContent>
    </S.Col>
    <Ticks numTicks={NUM_TICKS} startTime={0} endTime={duration} />

    <VerticalResizer nameColumnWidth={nameColumnWidth} onNameColumnWidthChange={onNameColumnWidthChange} />
  </S.HeaderRow>
);

export default Header;
