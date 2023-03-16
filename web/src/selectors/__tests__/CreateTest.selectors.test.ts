import Demo from 'models/Demo.model';
import {RootState} from 'redux/store';
import {ICreateTestState} from 'types/Test.types';
import {Plugins} from '../../constants/Plugins.constants';
import {ICreateTestStep} from '../../types/Plugins.types';
import CreateTestSelectors from '../CreateTest.selectors';

describe('CreateTestSelectors', () => {
  describe('selectStepList', () => {
    it('should return empty', () => {
      const result = CreateTestSelectors.selectStepList({
        createTest: {stepList: [] as ICreateTestStep[]} as ICreateTestState,
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });
  describe('selectPlugin', () => {
    it('should return pluginName', () => {
      const pluginName = Plugins.REST.name;

      const result = CreateTestSelectors.selectPlugin(
        {
          createTest: {pluginName} as ICreateTestState,
        } as RootState,
        [Demo()]
      );
      expect(result).toStrictEqual(Plugins.REST);
    });
  });
});
