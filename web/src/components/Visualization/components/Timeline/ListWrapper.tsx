import {FixedSizeList as List} from 'react-window';
import Header from './Header';
import SpanNodeFactory from './SpanNodeFactory';
import * as S from './Timeline.styled';
import {useTimeline} from './Timeline.provider';

const HEADER_HEIGHT = 242;

interface IProps {
  listRef: React.RefObject<List>;
}

const ListWrapper = ({listRef}: IProps) => {
  const {spans, viewEnd, viewStart} = useTimeline();

  return (
    <S.Container>
      <Header duration={viewEnd - viewStart} />
      <List
        height={window.innerHeight - HEADER_HEIGHT}
        itemCount={spans.length}
        itemData={spans}
        itemSize={32}
        ref={listRef}
        width="100%"
      >
        {SpanNodeFactory}
      </List>
    </S.Container>
  );
};

export default ListWrapper;
