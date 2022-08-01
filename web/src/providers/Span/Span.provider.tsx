import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useSelector} from 'react-redux';
import {useAppDispatch} from 'redux/hooks';
import {
  setAffectedSpans,
  setFocusedSpan,
  clearAffectedSpans,
  clearSelectedSpan,
  setSelectedSpan,
  setMatchedSpans,
  setSearchText,
} from 'redux/slices/Span.slice';
import SpanSelectors from 'selectors/Span.selectors';
import {TSpan} from 'types/Span.types';
import SpanService from '../../services/Span.service';
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
  onClearSelectedSpan(): void;
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
  onClearSelectedSpan: noop,
});

interface IProps {
  children: React.ReactNode;
}

export const useSpan = () => useContext(Context);

const SpanProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const {
    run: {trace: {spans = []} = {}},
  } = useTestRun();
  const selectedSpan = useSelector(SpanSelectors.selectSelectedSpan);
  const affectedSpans = useSelector(SpanSelectors.selectAffectedSpans);
  const focusedSpan = useSelector(SpanSelectors.selectFocusedSpan);
  const matchedSpans = useSelector(SpanSelectors.selectMatchedSpans);
  const searchText = useSelector(SpanSelectors.selectSearchText);

  useEffect(() => {
    return () => {
      dispatch(clearSelectedSpan());
      dispatch(clearAffectedSpans());
    };
  }, []);

  const onSelectSpan = useCallback(
    (spanId: string) => {
      const span = spans.find(({id}) => id === spanId);
      if (span) {
        dispatch(setSelectedSpan({span}));
      }
    },
    [dispatch, spans]
  );

  const onClearSelectedSpan = useCallback(() => {
    dispatch(clearSelectedSpan());
  }, [dispatch]);

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
      const spanIds = SpanService.searchSpanList(spans, query);
      dispatch(setMatchedSpans({spanIds}));
      dispatch(setSearchText({searchText: query}));
    },
    [dispatch, spans]
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
      onClearSelectedSpan,
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
      onClearSelectedSpan,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SpanProvider;
