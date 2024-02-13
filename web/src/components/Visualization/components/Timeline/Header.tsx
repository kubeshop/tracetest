import Ticks from './Ticks/Ticks';
import * as S from './TimelineV2.styled';

const NUM_TICKS = 5;

interface IProps {
  duration: number;
}

const Header = ({duration}: IProps) => (
  <S.HeaderRow>
    <S.Col>
      <S.HeaderContent>
        <S.HeaderTitle level={3}>Span</S.HeaderTitle>
      </S.HeaderContent>
      <S.Separator />
    </S.Col>
    <Ticks numTicks={NUM_TICKS} startTime={0} endTime={duration} />
  </S.HeaderRow>
);

export default Header;
