import {useEffect, useMemo, useRef} from 'react';
import {FixedSizeList as List} from 'react-window';
import AutoSizer, {Size} from 'react-virtualized-auto-sizer';

import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'models/AssertionResults.model';
import Header from './Header';
import ResultCard from './ResultCard';

interface IProps {
  onClose(): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onRevert(originalSelector: string): void;
  selectedSpan?: string;
  testSpec: TAssertionResultEntry;
}

const Content = ({
  onClose,
  onDelete,
  onEdit,
  onRevert,
  selectedSpan,
  testSpec,
  testSpec: {resultList, selector, spanIds},
}: IProps) => {
  const {
    isDeleted = false,
    isDraft = false,
    originalSelector = '',
    name = '',
  } = useAppSelector(state => TestSpecsSelectors.selectSpecBySelector(state, selector)) || {};
  const totalPassedChecks = useMemo(() => AssertionService.getTotalPassedChecks(resultList), [resultList]);
  const results = useMemo(() => Object.entries(AssertionService.getResultsHashedBySpanId(resultList)), [resultList]);

  const listRef = useRef<List>(null);

  useEffect(() => {
    if (listRef.current) {
      const index = results.findIndex(([spanId]) => spanId === selectedSpan);
      if (index !== -1) {
        listRef?.current?.scrollToItem(index, 'smart');
      }
    }
  }, [results, selectedSpan]);

  const itemSize = useMemo(() => {
    const [, checkResults = []] = results[0];

    return checkResults.length * 72.59 + 40 + 16;
  }, [results]);

  return (
    <>
      <Header
        affectedSpans={spanIds?.length ?? 0}
        assertionsFailed={totalPassedChecks?.false ?? 0}
        assertionsPassed={totalPassedChecks?.true ?? 0}
        isDeleted={isDeleted}
        isDraft={isDraft}
        onClose={onClose}
        onDelete={() => {
          onDelete(testSpec.selector);
          onClose();
        }}
        onEdit={() => {
          onEdit(testSpec, name);
        }}
        onRevert={() => {
          onRevert(originalSelector);
        }}
        selector={selector}
        title={!selector && !name ? 'All Spans' : name}
      />

      <AutoSizer>
        {({height, width}: Size) => (
          <List
            ref={listRef}
            height={height}
            itemCount={results.length}
            itemData={results}
            itemSize={itemSize}
            width={width}
          >
            {ResultCard}
          </List>
        )}
      </AutoSizer>
    </>
  );
};

export default Content;
