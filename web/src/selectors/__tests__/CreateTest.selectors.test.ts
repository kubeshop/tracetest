import {RootState} from '../../redux/store';
import {ICreateTestState} from '../../types/Plugins.types';
import CreateTestSelectors from '../CreateTest.selectors';

describe('CreateTestSelectors', () => {
  describe('selectStepList', () => {
    it('should return empty', () => {
      const result = CreateTestSelectors.selectStepList({
        createTest: {stepList: []} as ICreateTestState,
      } as RootState);
      expect(result).toStrictEqual([]);
    });
  });
  describe('selectPlugin', () => {
    it('should return pluginName', () => {
      const pluginName = 'pluginName';
      const result = CreateTestSelectors.selectPlugin({
        createTest: {pluginName} as ICreateTestState,
      } as RootState);
      expect(result).toStrictEqual(pluginName);
    });
  });
});
