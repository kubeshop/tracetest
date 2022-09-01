import {Tooltip} from 'antd';

import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import * as S from './Switch.styled';

interface IProps {
  onChange(type: VisualizationType): void;
  type: VisualizationType;
}

const Switch = ({onChange, type}: IProps) => (
  <S.Container>
    <Tooltip title="Graph view">
      <S.DAGIcon $isSelected={type === VisualizationType.Dag} onClick={() => onChange(VisualizationType.Dag)} />
    </Tooltip>
    <Tooltip title="Timeline view">
      <S.TimelineIcon
        $isSelected={type === VisualizationType.Timeline}
        onClick={() => onChange(VisualizationType.Timeline)}
      />
    </Tooltip>
  </S.Container>
);

export default Switch;
