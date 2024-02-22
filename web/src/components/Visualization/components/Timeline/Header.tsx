import {toPercent} from 'utils/Common';
import Ticks from './Ticks/Ticks';
import * as S from './Timeline.styled';
import VerticalResizer from './VerticalResizer/VerticalResizer';
import {useVerticalResizer} from './VerticalResizer.provider';

const NUM_TICKS = 5;

interface IProps {
  duration: number;
}

const Header = ({duration}: IProps) => {
  const {columnWidth} = useVerticalResizer();

  return (
    <S.HeaderRow style={{gridTemplateColumns: `${toPercent(columnWidth)} 1fr`}}>
      <S.Col>
        <S.HeaderContent>
          <S.HeaderTitle level={3}>Span</S.HeaderTitle>
        </S.HeaderContent>
      </S.Col>
      <Ticks numTicks={NUM_TICKS} startTime={0} endTime={duration} />

      <VerticalResizer />
    </S.HeaderRow>
  );
};

export default Header;
