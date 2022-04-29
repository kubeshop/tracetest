import LocalStorageGateway from '../gateways/LocalStorage.gateway';

export enum GuidedTours {
  Home = 'home',
  Assertion = 'assertion',
  Trace = 'trace',
  TestDetails = 'testDetails',
}

type TTour = Record<GuidedTours, boolean>;

const GUIDED_TOUR_KEY = 'guided_tour';

const {get, set} = LocalStorageGateway<TTour>(GUIDED_TOUR_KEY);

const defaultValue = Object.values(GuidedTours).reduce<TTour>(
  (acc, value) => ({
    ...acc,
    [value]: false,
  }),
  {} as TTour
);

const GuidedTourService = () => {
  return {
    get(): TTour {
      const guidedTour = get() || defaultValue;

      return guidedTour;
    },
    getIsComplete(tour: GuidedTours): boolean {
      const {[tour]: isComplete} = this.get();

      return isComplete;
    },
    save(tour: GuidedTours, isComplete = true): void {
      const guideTours = this.get();

      return set({
        ...guideTours,
        [tour]: isComplete,
      });
    },
    getSelectorStep(tour: GuidedTours, step: string, attribute = 'data-tour'): string {
      return `[${attribute}="${tour}_${step}"]`;
    },
    getStep(tour: GuidedTours, step: string): string {
      return `${tour}_${step}`;
    },
  };
};

export default GuidedTourService();
