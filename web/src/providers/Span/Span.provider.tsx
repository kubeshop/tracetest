import {noop} from 'lodash';
import {createContext, useCallback, useContext, useMemo} from 'react';
import {useSelector} from 'react-redux';
import {useAppDispatch} from 'redux/hooks';
import {setAffectedSpans, setFocusedSpan, clearAffectedSpans, setSelectedSpan} from 'redux/slices/Span.slice';
import SpanSelectors from 'selectors/Span.selectors';
import {TSpan} from 'types/Span.types';

interface IContext {
  selectedSpan?: TSpan;
  affectedSpans: string[];
  focusedSpan: string;
  onSelectSpan(spanId: TSpan): void;
  onSetAffectedSpans(spanIdList: string[]): void;
  onSetFocusedSpan(spanId: string): void;
  onClearAffectedSpans(): void;
}

export const Context = createContext<IContext>({
  affectedSpans: [],
  focusedSpan: '',
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
  const dispatch = useAppDispatch();
  const selectedSpan = useSelector(SpanSelectors.selectSelectedSpan);
  const affectedSpans = useSelector(SpanSelectors.selectAffectedSpans);
  const focusedSpan = useSelector(SpanSelectors.selectFocusedSpan);

  const onSelectSpan = useCallback(
    (span: TSpan) => {
      dispatch(setSelectedSpan({span}));
    },
    [dispatch]
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

  const value = useMemo<IContext>(
    () => ({
      selectedSpan,
      affectedSpans,
      focusedSpan,
      onSelectSpan,
      onSetAffectedSpans,
      onSetFocusedSpan,
      onClearAffectedSpans,
    }),
    [affectedSpans, focusedSpan, onClearAffectedSpans, onSelectSpan, onSetAffectedSpans, onSetFocusedSpan, selectedSpan]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
};

export default SpanProvider;
