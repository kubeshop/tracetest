import {ApartmentOutlined, ExpandOutlined, ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';
import {noop} from 'lodash';
import {useReactFlow} from 'react-flow-renderer';
import AffectedSpanControls from './AffectedSpanControls';
import * as S from './DAG.styled';
import {ControlsMode} from './DAG.styled';

interface IProps {
  isMiniMapActive?: boolean;
  mode?: ControlsMode;
  onMiniMapToggle?: () => void;
}

const Controls = ({mode = 'dag', isMiniMapActive = false, onMiniMapToggle = noop}: IProps) => {
  const {fitView, zoomIn, zoomOut} = useReactFlow();
  return (
    <>
      <S.SelectorControls mode={mode}>
        <AffectedSpanControls />
      </S.SelectorControls>
      {mode === 'dag' && (
        <S.Controls>
          <S.ZoomButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
          <S.ZoomButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
          <S.ZoomButton icon={<ExpandOutlined />} onClick={() => fitView()} type="text" />
          <S.ZoomButton
            icon={<ApartmentOutlined />}
            onClick={onMiniMapToggle}
            type="text"
            $isActive={isMiniMapActive}
          />
        </S.Controls>
      )}
    </>
  );
};

export default Controls;
