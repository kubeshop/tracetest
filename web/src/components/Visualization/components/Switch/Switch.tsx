import {Tooltip} from 'antd';
import {VisualizationType} from 'components/RunDetailTrace/RunDetailTrace';
import {MAX_DAG_NODES} from 'constants/Visualization.constants';
import * as S from './Switch.styled';

interface IProps {
  isDAGDisabled: boolean;
  onChange(type: VisualizationType): void;
  type: VisualizationType;
  totalSpans?: number;
}

const Switch = ({isDAGDisabled, onChange, type, totalSpans = 0}: IProps) => (
  <S.Container>
    <Tooltip
      title={
        isDAGDisabled
          ? `The Graph view has a limit of ${MAX_DAG_NODES} spans. Your current trace has ${totalSpans} spans.`
          : 'Graph view'
      }
      placement="right"
    >
      <S.DAGIcon
        $isDisabled={isDAGDisabled}
        $isSelected={type === VisualizationType.Dag}
        onClick={() => !isDAGDisabled && onChange(VisualizationType.Dag)}
      />
    </Tooltip>
    <Tooltip title="Timeline view" placement="right">
      <S.TimelineIcon
        $isSelected={type === VisualizationType.Timeline}
        onClick={() => onChange(VisualizationType.Timeline)}
      />
    </Tooltip>
  </S.Container>
);

export default Switch;
