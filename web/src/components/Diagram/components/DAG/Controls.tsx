import {ApartmentOutlined, ExpandOutlined, ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';
import {useReactFlow} from 'react-flow-renderer';
import * as S from './DAG.styled';
import AffectedSpanControls from './AffectedSpanControls';

interface IProps {
  isMiniMapActive: boolean;
  onMiniMapToggle(): void;
}

const Controls = ({isMiniMapActive, onMiniMapToggle}: IProps) => {
  const {fitView, zoomIn, zoomOut} = useReactFlow();

  return (
    <>
      <S.SelectorControls>
        <AffectedSpanControls />
      </S.SelectorControls>
      <S.Controls>
        <S.ZoomButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
        <S.ZoomButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
        <S.ZoomButton icon={<ExpandOutlined />} onClick={() => fitView()} type="text" />
        <S.ZoomButton icon={<ApartmentOutlined />} onClick={onMiniMapToggle} type="text" $isActive={isMiniMapActive} />
      </S.Controls>
    </>
  );
};

export default Controls;
