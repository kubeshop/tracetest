import {TourProvider} from '@reactour/tour';
import {doArrow} from './doArrow';
import './index.css';
import {NextButton} from './NextButton';
import {PreviousButton} from './PreviousButton';

export const GuidedTourProvider: React.FC = ({children}) => {
  return (
    <TourProvider
      steps={[]}
      maskClassName="tour-mask"
      prevButton={PreviousButton}
      nextButton={NextButton}
      showCloseButton={false}
      onClickClose={clickProps => clickProps.setCurrentStep(0)}
      styles={{
        badge: props => ({...props, display: 'none'}),
        navigation: props => ({...props, display: 'none', margin: 0}),
        controls: props => ({...props, margin: 0, padding: 16, paddingTop: 0}),
        popover: (base, state) => {
          return {
            ...base,
            padding: 0,
            borderRadius: 10,
            ...doArrow(state?.position, state?.verticalAlign, state?.horizontalAlign),
          };
        },
      }}
    >
      {children}
    </TourProvider>
  );
};
