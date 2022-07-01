import {LeftOutlined, RightOutlined} from '@ant-design/icons';
import {useSpan} from 'providers/Span/Span.provider';
import * as S from './DAG.styled';
import useControls from './hooks/useControls';

const AffectedSpanControls = () => {
  const {affectedSpans, onSelectSpan} = useSpan();
  const {handleNextSpan, handlePrevSpan, indexOfFocused} = useControls({onSelectSpan});

  return affectedSpans.length ? (
    <div>
      <S.ToggleButton onClick={handlePrevSpan} icon={<LeftOutlined />} type="text" />
      <S.ToggleButton onClick={handleNextSpan} icon={<RightOutlined />} type="text" />
      <S.FocusedText>
        {indexOfFocused} of {affectedSpans.length} total
      </S.FocusedText>
    </div>
  ) : null;
};

export default AffectedSpanControls;
