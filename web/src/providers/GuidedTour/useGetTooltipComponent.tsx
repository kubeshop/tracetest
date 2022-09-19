import React, {useCallback, useMemo} from 'react';
import {TooltipRenderProps} from 'react-joyride';
import {ReactJoyrideTooltip} from '../../components/GuidedTour/StepContent';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';

export function useGetTooltipComponent(tour: GuidedTours): React.FC<TooltipRenderProps> {
  const onComplete = useCallback(() => GuidedTourService.save(tour), [tour]);
  return useMemo(() => ReactJoyrideTooltip(onComplete), [onComplete]);
}
