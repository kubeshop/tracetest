import {ParentSize} from '@visx/responsive';

import {TSpan} from 'types/Span.types';
import * as S from './Timeline.styled';
import Visualization from './Visualization';
import Navigation from '../Navigation';

export interface IProps {
  isMatchedMode: boolean;
  matchedSpans: string[];
  onNavigateToSpan(spanId: string): void;
  onNodeClick(spanId: string): void;
  selectedSpan: string;
  spans: TSpan[];
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
