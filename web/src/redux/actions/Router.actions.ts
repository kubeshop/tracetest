import {createAsyncThunk} from '@reduxjs/toolkit';
import {parse, ParsedQuery, stringify} from 'query-string';
import {Params} from 'react-router-dom';
import {push} from 'redux-first-history';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {RouterSearchFields} from '../../constants/Common.constants';
import SpanSelectors from '../../selectors/Span.selectors';
import {decryptString} from '../../utils/Common';
import {setSelectedSpan} from '../slices/Span.slice';
import {setSelectedAssertion} from '../slices/TestDefinition.slice';
import {RootState} from '../store';

export interface IQuery {
  search: ParsedQuery<string>;
  params: Readonly<Params<string>>;
}

const RouterActions = () => ({
  updateSelectedAssertion: createAsyncThunk<void, IQuery>(
    'router/addAssertionResult',
    async ({search}, {getState, dispatch}) => {
      const {[RouterSearchFields.SelectedAssertion]: selector = ''} = search;

      const decryptedSelector = decryptString(String(selector));

      const assertionResult = TestDefinitionSelectors.selectAssertionBySelector(
        getState() as RootState,
        decryptedSelector
      );

      if (assertionResult) dispatch(setSelectedAssertion(assertionResult));
      else if (!selector) dispatch(setSelectedAssertion());
    }
  ),
  updateSelectedSpan: createAsyncThunk<void, IQuery>(
    'router/updateSelectedSpan',
    async ({params: {testId = '', runId = ''}, search}, {getState, dispatch}) => {
      const {[RouterSearchFields.SelectedSpan]: spanId = ''} = search;
      const state = getState() as RootState;
      const span = SpanSelectors.selectSpanById(state, String(spanId), testId, runId);

      if (span) dispatch(setSelectedSpan({span}));
    }
  ),
  updateSearch: createAsyncThunk<void, Partial<Record<RouterSearchFields, any>>>(
    'router/updateSearch',
    async (newSearch, {dispatch, getState}) => {
      const {
        router: {location},
      } = getState() as RootState;

      const search = parse(location?.search || '');

      const filteredSearch = Object.entries({
        ...search,
        ...newSearch,
      }).reduce((acc, [key, value]) => (value ? {...acc, [key]: value} : acc), {});

      await dispatch(
        push({
          search: stringify(filteredSearch),
        })
      );
    }
  ),
});

export default RouterActions();
