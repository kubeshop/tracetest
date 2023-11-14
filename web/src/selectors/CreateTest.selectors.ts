import {createSelector} from '@reduxjs/toolkit';
import {getDemoByPluginMap} from 'constants/Demo.constants';
import {Plugins} from 'constants/Plugins.constants';
import Demo from 'models/Demo.model';
import {RootState} from 'redux/store';
import {IPlugin} from 'types/Plugins.types';

const createTestSelectorsStateSelector = (state: RootState) => state.createTest;

const selectDemos = (state: RootState, demos: Demo[]) => demos;

const CreateTestSelectors = () => ({
  selectPlugin: createSelector(createTestSelectorsStateSelector, selectDemos, ({pluginName}, demos) => {
    const demoByPluginMap = getDemoByPluginMap(demos);
    const demoList = demoByPluginMap[pluginName];
    return {...Plugins[pluginName], demoList} as IPlugin;
  }),
  selectDraftTest: createSelector(createTestSelectorsStateSelector, ({draftTest}) => draftTest),
  selectIsFormValid: createSelector(createTestSelectorsStateSelector, ({isFormValid}) => isFormValid),
});

export default CreateTestSelectors();
