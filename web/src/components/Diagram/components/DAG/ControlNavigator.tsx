import {LeftOutlined, RightOutlined} from '@ant-design/icons';

import {useSpan} from 'providers/Span/Span.provider';
import * as S from './DAG.styled';
import useControls from './hooks/useControls';
import {useListenToArrowKeysDown} from './useListenToArrowKeysDown';

const ControlNavigator = () => {
  const {affectedSpans} = useSpan();
  const {handleNextSpan, handlePrevSpan, indexOfFocused} = useControls();
  useListenToArrowKeysDown(affectedSpans, handleNextSpan, handlePrevSpan);

  if (!affectedSpans.length) return <div />;

  return (
    <S.NavigatorPanel>
      <S.ToggleButton id="span-back" onClick={handlePrevSpan} icon={<LeftOutlined />} type="text" />
      <S.ToggleButton id="span-forward" onClick={handleNextSpan} icon={<RightOutlined />} type="text" />
      <S.FocusedText strong>
        {indexOfFocused} of {affectedSpans.length} total
      </S.FocusedText>
    </S.NavigatorPanel>
  );
};

export default ControlNavigator;
