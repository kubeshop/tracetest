import {noop} from 'lodash';
import {createContext, useCallback, useContext, useEffect, useMemo} from 'react';
import {useSelector} from 'react-redux';
import {useAppDispatch} from 'redux/hooks';
import {setMatchedSpans, setFocusedSpan, clearMatchedSpans, clearSelectedSpan} from 'redux/slices/Span.slice';
import SpanSelectors from 'selectors/Span.selectors';
import {TSpan} from 'types/Span.types';
import {RouterSearchFields} from '../../constants/Common.constants';
import RouterActions from '../../redux/actions/Router.actions';

interface IContext {
  selectedSpan?: TSpan;
  matchedSpans: string[];
  focusedSpan: string;
  onSelectSpan(spanId: string): void;
  onSetMatchedSpans(spanIdList: string[]): void;
  onSetFocusedSpan(spanId: string): void;
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
    }),
    [
      matchedSpans,
      focusedSpan,
      onClearMatchedSpans,
      onSelectSpan,
      onSetMatchedSpans,
      onSetFocusedSpan,
      selectedSpan,
      onClearSelectedSpan,
    ]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SpanProvider;
