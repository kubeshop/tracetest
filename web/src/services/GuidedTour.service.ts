import {GuidedTourByPathname, GuidedTours} from 'constants/GuidedTour';

const GuidedTourService = {
  getByPathName: (pathname: string): GuidedTours | undefined => {
    const [, value] =
      Object.entries(GuidedTourByPathname).find(([key]) => {
        const regex = new RegExp(key);
        return regex.test(pathname);
      }) || [];
    return value;
  },

  getStepSelector(step: string, attribute = 'data-tour'): string {
    return `[${attribute}="${step}"]`;
  },
};

export default GuidedTourService;
