import {Tooltip} from 'antd';

import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import ExperimentalFeature from '../../../../utils/ExperimentalFeature';
import * as S from './Switch.styled';

interface IProps {
  type: VisualizationType;
  onChange(type: VisualizationType): void;
}

const Switch = ({onChange, type}: IProps) => (
  <S.Container>
    <Tooltip title="Graph view" placement="right">
      <S.DAGIcon $isSelected={type === VisualizationType.Dag} onClick={() => onChange(VisualizationType.Dag)} />
    </Tooltip>
    <Tooltip title="Timeline view" placement="right">
      <S.TimelineIcon
        $isSelected={type === VisualizationType.Timeline}
        onClick={() => onChange(VisualizationType.Timeline)}
      />
    </Tooltip>
    {ExperimentalFeature.isEnabled('flamegraph') && (
      <Tooltip title="Flame view" placement="right">
        <S.FlameIcon $isSelected={type === VisualizationType.Flame} onClick={() => onChange(VisualizationType.Flame)} />
      </Tooltip>
    )}
  </S.Container>
);

export default Switch;
