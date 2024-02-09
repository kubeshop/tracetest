import Navigation from '../Navigation';
import {useTimeline} from './Timeline.provider';

const NavigationWrapper = () => {
  const {matchedSpans, onSpanNavigation, selectedSpan} = useTimeline();

  return <Navigation matchedSpans={matchedSpans} onNavigateToSpan={onSpanNavigation} selectedSpan={selectedSpan} />;
};

export default NavigationWrapper;
