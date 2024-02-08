import Navigation from '../Navigation';
import {useTimeline} from './Timeline.provider';

const NavigationWrapper = () => {
  const {matchedSpans, onNavigateToSpan, selectedSpan} = useTimeline();

  return <Navigation matchedSpans={matchedSpans} onNavigateToSpan={onNavigateToSpan} selectedSpan={selectedSpan} />;
};

export default NavigationWrapper;
