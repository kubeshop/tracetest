import Demo from 'models/Demo.model';
import {RootState} from 'redux/store';
import {ICreateTestState} from 'types/Test.types';
import {Plugins} from '../../constants/Plugins.constants';
import CreateTestSelectors from '../CreateTest.selectors';

describe('CreateTestSelectors', () => {
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
