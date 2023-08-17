import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useSelector} from 'react-redux';
import {useAppDispatch} from 'redux/hooks';
import {setMatchedSpans, setFocusedSpan, clearMatchedSpans, clearSelectedSpan} from 'redux/slices/Span.slice';
import SpanSelectors from 'selectors/Span.selectors';
import {RouterSearchFields} from 'constants/Common.constants';
import Span from 'models/Span.model';
import RouterActions from 'redux/actions/Router.actions';
import TracetestAPI from 'redux/apis/Tracetest';
import SelectedSpans from 'models/SelectedSpans.model';

const {useLazyGetSelectedSpansQuery} = TracetestAPI.instance;

interface IContext {
  selectedSpan?: Span;
  triggerSelectorResult?: SelectedSpans;
  matchedSpans: string[];
  focusedSpan: string;
  isTriggerSelectorError: boolean;
  isTriggerSelectorLoading: boolean;
  onSelectSpan(spanId: string): void;
  onSetMatchedSpans(spanIdList: string[]): void;
  onSetFocusedSpan(spanId: string): void;
  onTriggerSelector(query: string, testId: string, runId: string): void;
  onClearMatchedSpans(): void;
  onClearSelectedSpan(): void;
}

export const Context = createContext<IContext>({
  matchedSpans: [],
  focusedSpan: '',
  onSelectSpan: noop,
  onSetFocusedSpan: noop,
  onClearMatchedSpans: noop,
  onSetMatchedSpans: noop,
  onClearSelectedSpan: noop,
  onTriggerSelector: noop,
  isTriggerSelectorError: false,
  isTriggerSelectorLoading: false,
});

interface IProps {
  children: React.ReactNode;
}

export const useSpan = (): IContext => useContext(Context);

const SpanProvider = ({children}: IProps) => {
  const dispatch = useAppDispatch();
  const selectedSpan = useSelector(SpanSelectors.selectSelectedSpan);
  const matchedSpans = useSelector(SpanSelectors.selectMatchedSpans);
  const focusedSpan = useSelector(SpanSelectors.selectFocusedSpan);
  const [
    triggerSelector,
    {data: triggerSelectorResult, isError: isTriggerSelectorError, isFetching: isTriggerSelectorLoading},
  ] = useLazyGetSelectedSpansQuery();

  useEffect(() => {
    return () => {
      dispatch(RouterActions.updateSearch({[RouterSearchFields.SelectedSpan]: ''}));
      dispatch(clearMatchedSpans());
    };
  }, [dispatch]);

  const onSelectSpan = useCallback(
    (spanId: string) => {
      dispatch(RouterActions.updateSearch({[RouterSearchFields.SelectedSpan]: spanId}));
    },
    [dispatch]
  );

  const onClearSelectedSpan = useCallback(() => {
    dispatch(clearSelectedSpan());
  }, [dispatch]);

  const onSetMatchedSpans = useCallback(
    (spanIds: string[]) => {
      dispatch(setMatchedSpans({spanIds}));
    },
    [dispatch]
  );

  const onTriggerSelector = useCallback(
    async (query: string, testId: string, runId: string) => {
      const {spanIds = []} = await triggerSelector({
        query,
        testId,
        runId,
      }).unwrap();
      onSetMatchedSpans(spanIds);
    },
    [onSetMatchedSpans, triggerSelector]
  );

  const onSetFocusedSpan = useCallback(
    (spanId: string) => {
      dispatch(setFocusedSpan({spanId}));
    },
    [dispatch]
  );

  const onClearMatchedSpans = useCallback(() => {
    dispatch(clearMatchedSpans());
  }, [dispatch]);

  const value = useMemo<IContext>(
    () => ({
      selectedSpan,
      matchedSpans,
      focusedSpan,
      onSelectSpan,
      onSetMatchedSpans,
      onSetFocusedSpan,
      onClearMatchedSpans,
      onClearSelectedSpan,
      onTriggerSelector,
      triggerSelectorResult,
      isTriggerSelectorError,
      isTriggerSelectorLoading,
    }),
    [
      selectedSpan,
      matchedSpans,
      focusedSpan,
      onSelectSpan,
      onSetMatchedSpans,
      onSetFocusedSpan,
      onClearMatchedSpans,
      onClearSelectedSpan,
      onTriggerSelector,
      triggerSelectorResult,
      isTriggerSelectorError,
      isTriggerSelectorLoading,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SpanProvider;
