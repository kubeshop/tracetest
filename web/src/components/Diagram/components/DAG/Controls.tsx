import {ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';
import {useReactFlow} from 'react-flow-renderer';
import * as S from './DAG.styled';
import AffectedSpanControls from './AffectedSpanControls';

const Controls = () => {
  const {zoomIn, zoomOut} = useReactFlow();

  return (
    <>
      <S.SelectorControls>
        <AffectedSpanControls />
      </S.SelectorControls>
      <S.Controls>
        <S.ZoomButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
        <S.ZoomButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
      </S.Controls>
    </>
  );
};

export default Controls;
