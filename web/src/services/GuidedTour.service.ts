import {useLocation} from 'react-router-dom';
import LocalStorageGateway from '../gateways/LocalStorage.gateway';

export enum GuidedTours {
  Home = 'home',
  Environment = 'environment',
  Trace = 'trace',
  TestDetails = 'testDetails',
}

export const GuidedToursPathnameMap = {
  '/run/': GuidedTours.Trace,
  '/test/': GuidedTours.TestDetails,
  '/': GuidedTours.Home,
};

type TTour = Record<GuidedTours, boolean>;

const GUIDED_TOUR_KEY = 'guided_tour';

const {get, set} = LocalStorageGateway<TTour>(GUIDED_TOUR_KEY);

export const defaultValue = Object.values(GuidedTours).reduce<TTour>(
  (acc, value) => ({
    ...acc,
    [value]: false,
  }),
  {} as TTour
);

const GuidedTourService = () => {
  return {
    getByPathName: (pathname: string): GuidedTours => {
      const [, value = GuidedTours.Home] =
        Object.entries(GuidedToursPathnameMap).find(([key]) => pathname.includes(key)) || [];
      return value;
    },
    get(): TTour {
      return get() || defaultValue;
    },
    useGetCurrentOnboardingLocation(): GuidedTours {
      const location = useLocation();
      return this.getByPathName(location.pathname);
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
