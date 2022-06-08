import GuidedTourService, {defaultValue, GuidedTours} from '../GuidedTour.service';

describe('GuidedTourService', () => {
  it('should get the default guided tour object from the local storage', () => {
    const guidedTour = GuidedTourService.get();

    expect(guidedTour).toEqual(defaultValue);
  });

  it('should return the updated value after save', () => {
    GuidedTourService.save(GuidedTours.Home);
    const guidedTour = GuidedTourService.get();

    expect(guidedTour.home).toBeTruthy();
  });

  it('getIsComplete should return after the tour is completed', () => {
    expect(GuidedTourService.getIsComplete(GuidedTours.Trace)).toBeFalsy();
    GuidedTourService.save(GuidedTours.Trace);

    expect(GuidedTourService.getIsComplete(GuidedTours.Trace)).toBeTruthy();
  });

  it('getSelectorStep should return the identifier for a specific guided tour step', () => {
    const selector = GuidedTourService.getSelectorStep(GuidedTours.Home, 'step1');

    expect(selector).toEqual('[data-tour="home_step1"]');
  });

  it('getStep should return the step identifier for a specific guided tour', () => {
    const step = GuidedTourService.getStep(GuidedTours.Home, 'step1');

    expect(step).toEqual('home_step1');
  });
});
