import {createSelector} from '@reduxjs/toolkit';
import {getDemoByPluginMap} from 'constants/Demo.constants';
import {Plugins} from 'constants/Plugins.constants';
import Demo from 'models/Demo.model';
import {RootState} from 'redux/store';
import {IPlugin} from 'types/Plugins.types';

const createTestSelectorsStateSelector = (state: RootState) => state.createTest;

const selectDemo = (state: RootState, demo: Demo) => demo;

const CreateTestSelectors = () => ({
  selectStepList: createSelector(createTestSelectorsStateSelector, ({stepList}) => stepList),
  selectPlugin: createSelector(createTestSelectorsStateSelector, selectDemo, ({pluginName}, demo) => {
    const demoByPluginMap = getDemoByPluginMap(demo);
    const demoList = demoByPluginMap[pluginName];
    return {...Plugins[pluginName], demoList} as IPlugin;
  }),
  selectStepNumber: createSelector(createTestSelectorsStateSelector, ({stepNumber}) => stepNumber),
  selectDraftTest: createSelector(createTestSelectorsStateSelector, ({draftTest}) => draftTest),
  selectIsFormValid: createSelector(createTestSelectorsStateSelector, ({isFormValid}) => isFormValid),
  selectActiveStep: createSelector(
    createTestSelectorsStateSelector,
    ({stepList, stepNumber}) => stepList[stepNumber]?.id || ''
  ),
});

export default CreateTestSelectors();
