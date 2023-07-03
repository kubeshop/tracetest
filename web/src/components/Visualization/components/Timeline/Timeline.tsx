import {ParentSize} from '@visx/responsive';
import {NodeTypesEnum} from 'constants/Visualization.constants';
import Span from 'models/Span.model';
import * as S from './Timeline.styled';
import Visualization from './Visualization';
import Navigation from '../Navigation';

export interface IProps {
  isMatchedMode: boolean;
  matchedSpans: string[];
  nodeType: NodeTypesEnum;
  onNavigateToSpan(spanId: string): void;
  onNodeClick(spanId: string): void;
  selectedSpan: string;
  spans: Span[];
  width?: number;
}

const Timeline = (props: IProps) => {
  const {isMatchedMode, matchedSpans, onNavigateToSpan, selectedSpan} = props;

  return (
    <S.Container $showMatched={isMatchedMode}>
      <Navigation matchedSpans={matchedSpans} onNavigateToSpan={onNavigateToSpan} selectedSpan={selectedSpan} />
      <ParentSize parentSizeStyles={{height: '100%', overflowY: 'scroll', paddingTop: 32, width: '100%'}}>
        {({width}) => <Visualization {...props} width={width} />}
      </ParentSize>
    </S.Container>
  );
};

export default Timeline;
