import {ApartmentOutlined, ExpandOutlined, ZoomInOutlined, ZoomOutOutlined} from '@ant-design/icons';
import {noop} from 'lodash';
import {useReactFlow} from 'react-flow-renderer';

import ControlNavigator from './ControlNavigator';
import * as S from './DAG.styled';

interface IProps {
  isMiniMapActive?: boolean;
  mode?: 'timeline' | 'dag';
  onMiniMapToggle?: () => void;
}

const Controls = ({isMiniMapActive = false, mode = 'dag', onMiniMapToggle = noop}: IProps) => {
  const {fitView, zoomIn, zoomOut} = useReactFlow();

  return (
    <>
      <ControlNavigator />

      {mode === 'dag' && (
        <S.DAGActionsPanel>
          <S.ActionButton icon={<ZoomInOutlined />} onClick={() => zoomIn()} type="text" />
          <S.ActionButton icon={<ZoomOutOutlined />} onClick={() => zoomOut()} type="text" />
          <S.ActionButton icon={<ExpandOutlined />} onClick={() => fitView()} type="text" />
          <S.ActionButton
            icon={<ApartmentOutlined />}
            onClick={onMiniMapToggle}
            type="text"
            $isActive={isMiniMapActive}
          />
        </S.DAGActionsPanel>
      )}
    </>
  );
};

export default Controls;
