import {LeftOutlined, RightOutlined, ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';
import {useReactFlow} from 'react-flow-renderer';
import {useSpan} from 'providers/Span/Span.provider';
import * as S from './DAG.styled';
import useControls from './hooks/useControls';

interface IProps {
  onSelectSpan(spanId: string): void;
}

const Controls = ({onSelectSpan}: IProps) => {
  const {zoomIn, zoomOut} = useReactFlow();
  const {affectedSpans} = useSpan();
  const {handleNextSpan, handlePrevSpan, indexOfFocused} = useControls({onSelectSpan});

  return (
    <>
      {Boolean(affectedSpans.length) && (
        <S.SelectorControls>
          <S.ToggleButton onClick={handlePrevSpan} icon={<LeftOutlined />} type="text" />
          <S.ToggleButton onClick={handleNextSpan} icon={<RightOutlined />} type="text" />
          <S.FocusedText>
            {indexOfFocused} of {affectedSpans.length}
          </S.FocusedText>
        </S.SelectorControls>
      )}
      <S.Controls>
        <S.ZoomButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
        <S.ZoomButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
      </S.Controls>
    </>
  );
};

export default Controls;
