import {Link} from 'react-router-dom';
import HomeAnalyticsService from 'services/Analytics/HomeAnalytics.service';
import {Steps} from '../../components/GuidedTour/homeStepList';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import * as S from './Home.styled';

const {onCreateTestClick} = HomeAnalyticsService;

const HomeActions = () => {
  return (
    <S.ActionContainer>
      <Link to="/test/create">
        <S.CreateTestButton
          data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.CreateTest)}
          data-cy="create-test-button"
          type="primary"
          onClick={() => {
            onCreateTestClick();
          }}
        >
          Create Test
        </S.CreateTestButton>
      </Link>
    </S.ActionContainer>
  );
};

export default HomeActions;
