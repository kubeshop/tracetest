import {useSpan} from 'providers/Span/Span.provider';
import {TPanel, TPanelComponentProps} from '../ResizablePanels/ResizablePanels';
import SpanDetail from '../SpanDetail/SpanDetail';
import * as S from '../RunDetailTrace/RunDetailTrace.styled';

const SpanDetailsPanel = ({size: {isOpen}}: TPanelComponentProps) => {
  const {selectedSpan} = useSpan();

  return (
    <S.PanelContainer $isOpen={isOpen}>
      <SpanDetail span={selectedSpan} />
    </S.PanelContainer>
  );
};

export const getSpanDetailsPanel = (): TPanel => ({
  name: 'SPAN_DETAILS',
  minSize: 15,
  maxSize: 320,
  position: 'left',
  component: props => <SpanDetailsPanel {...props} />,
});

export default SpanDetailsPanel;
