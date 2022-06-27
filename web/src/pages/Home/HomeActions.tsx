import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {Steps} from '../../components/GuidedTour/homeStepList';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import * as S from './Home.styled';

const {onCreateTestClick} = HomeAnalyticsService;

const HomeActions = () => {
  return (
    <S.ActionContainer>
      <S.CreateTestButton
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.CreateTest)}
        data-cy="create-test-button"
        type="primary"
        href="/create-test"
        onClick={() => {
          onCreateTestClick();
        }}
      >
        Create Test
      </S.CreateTestButton>
    </S.ActionContainer>
  );
};

export default HomeActions;
