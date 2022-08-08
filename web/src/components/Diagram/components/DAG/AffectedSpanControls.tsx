import {LeftOutlined, RightOutlined} from '@ant-design/icons';
import {useSpan} from 'providers/Span/Span.provider';
import * as S from './DAG.styled';
import useControls from './hooks/useControls';
import {useListenToArrowKeysDown} from './useListenToArrowKeysDown';

const AffectedSpanControls = () => {
  const {affectedSpans, onSelectSpan} = useSpan();
  const {handleNextSpan, handlePrevSpan, indexOfFocused} = useControls({onSelectSpan});
  useListenToArrowKeysDown(affectedSpans, handleNextSpan, handlePrevSpan);
  return (
    <div>
      <S.ToggleButton id="span-back" onClick={handlePrevSpan} icon={<LeftOutlined />} type="text" />
      <S.ToggleButton id="span-forward" onClick={handleNextSpan} icon={<RightOutlined />} type="text" />
      <S.FocusedText>
        {indexOfFocused} of {affectedSpans.length} total
      </S.FocusedText>
    </div>
  );
};

export default AffectedSpanControls;
