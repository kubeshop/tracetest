import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo, useState} from 'react';
import {useStoreActions} from 'react-flow-renderer';
import {useSelector} from 'react-redux';
import {useAppDispatch} from 'redux/hooks';
import {setAffectedSpans, setFocusedSpan, clearAffectedSpans, setSelectedSpan} from 'redux/slices/Span.slice';
import SpanSelectors from 'selectors/Span.selectors';
import {TSpan} from 'types/Span.types';
import SpanService from '../../services/Span.service';
// import {useLazySearchSpansQuery} from '../../redux/apis/TraceTest.api';
// import {useTestDefinition} from '../TestDefinition/TestDefinition.provider';
import {useTestRun} from '../TestRun/TestRun.provider';

interface IContext {
  selectedSpan?: TSpan;
  affectedSpans: string[];
  focusedSpan: string;
  matchedSpans: string[];
  searchText: string;
  onSearch(searchText: string): void;
  onSelectSpan(spanId: string): void;
  onSetAffectedSpans(spanIdList: string[]): void;
  onSetFocusedSpan(spanId: string): void;
  onClearAffectedSpans(): void;
}

export const Context = createContext<IContext>({
  affectedSpans: [],
  focusedSpan: '',
  matchedSpans: [],
  searchText: '',
  onSearch: noop,
  onSelectSpan: noop,
  onSetFocusedSpan: noop,
  onClearAffectedSpans: noop,
  onSetAffectedSpans: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useSpan = () => useContext(Context);

const SpanProvider = ({children}: IProps) => {
  const [searchText, setSearchText] = useState('');
  const [matchedSpans, setMatchedSpans] = useState<string[]>([]);

  const dispatch = useAppDispatch();
  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const {
    run: {trace: {spans = []} = {}},
  } = useTestRun();
  // const {test: {id: testId = ''} = {}} = useTestDefinition();
  const selectedSpan = useSelector(SpanSelectors.selectSelectedSpan);
  const affectedSpans = useSelector(SpanSelectors.selectAffectedSpans);
  const focusedSpan = useSelector(SpanSelectors.selectFocusedSpan);
  // const [triggerSearch, {data: foundSpans = []}] = useLazySearchSpansQuery();

  const onSelectSpan = useCallback(
    (spanId: string) => {
      const span = spans.find(({id}) => id === spanId);
      if (span) {
        addSelected([{id: span?.id}]);
        dispatch(setSelectedSpan({span}));
      }
    },
    [addSelected, dispatch, spans]
  );

  const onSetAffectedSpans = useCallback(
    (spanIds: string[]) => {
      dispatch(setAffectedSpans({spanIds}));
    },
    [dispatch]
  );

  const onSetFocusedSpan = useCallback(
    (spanId: string) => {
      dispatch(setFocusedSpan({spanId}));
    },
    [dispatch]
  );

  const onClearAffectedSpans = useCallback(() => {
    dispatch(clearAffectedSpans());
  }, [dispatch]);

  const onSearch = useCallback(
    (query: string) => {
      const spanIdList = SpanService.searchSpanList(spans, query);
      setMatchedSpans(spanIdList);
      setSearchText(query);

      console.log('@@query', query, spanIdList);
    },
    [spans]
  );

  const value = useMemo<IContext>(
    () => ({
      selectedSpan,
      affectedSpans,
      matchedSpans,
      focusedSpan,
      searchText,
      onSearch,
      onSelectSpan,
      onSetAffectedSpans,
      onSetFocusedSpan,
      onClearAffectedSpans,
    }),
    [
      affectedSpans,
      focusedSpan,
      matchedSpans,
      onClearAffectedSpans,
      onSearch,
      onSelectSpan,
      onSetAffectedSpans,
      onSetFocusedSpan,
      searchText,
      selectedSpan,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SpanProvider;
