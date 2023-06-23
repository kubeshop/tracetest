import {filter, findLastIndex, flow} from 'lodash';
import {TestRunStage, TraceEventType} from 'constants/TestRunEvents.constants';
import AssertionResults from 'models/AssertionResults.model';
import LinterResult from 'models/LinterResult.model';
import {isRunStateAnalyzingError, isRunStateStopped, isRunStateSucceeded} from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import TestRunOutput from 'models/TestRunOutput.model';
import {TAnalyzerErrorsBySpan, TTestOutputsBySpan, TTestRunState, TTestSpecsBySpan} from 'types/TestRun.types';

const TestRunService = () => ({
  shouldDisplayTraceEvents(state: TTestRunState, numberOfSpans: number) {
    const isStateSucceededOrStopped =
      isRunStateSucceeded(state) || isRunStateStopped(state) || isRunStateAnalyzingError(state);
    return !isStateSucceededOrStopped || (!numberOfSpans && isStateSucceededOrStopped);
  },

  getTestRunEventsByStage(events: TestRunEvent[], stage: TestRunStage) {
    return events.filter(event => event.stage === stage);
  },

  getTestRunTraceEvents(events: TestRunEvent[]): TestRunEvent[] {
    return flow([this.getTestRunEventsWithoutFetching, this.getTestRunEventsWithLastPolling])(events);
  },

  getTestRunEventsWithoutFetching(events: TestRunEvent[]) {
    return filter(
      events,
      event =>
        !(
          [TraceEventType.FETCHING_START, TraceEventType.FETCHING_ERROR, TraceEventType.FETCHING_SUCCESS] as string[]
        ).includes(event.type)
    );
  },

  getTestRunEventsWithLastPolling(events: TestRunEvent[]) {
    const lastPollingIndex = findLastIndex(events, event => event.type === TraceEventType.POLLING_ITERATION_INFO);
    if (lastPollingIndex === -1) return events;

    const eventsWithoutPolling = filter(events, event => event.type !== TraceEventType.POLLING_ITERATION_INFO);
    const newIndex = lastPollingIndex - (events.length - eventsWithoutPolling.length) + 1;
    eventsWithoutPolling.splice(newIndex, 0, events[lastPollingIndex]);

    return eventsWithoutPolling;
  },

  getAnalyzerErrorsHashedBySpan(linterResult: LinterResult): TAnalyzerErrorsBySpan {
    return linterResult.plugins
      .flatMap(plugin => plugin.rules.map(rule => ({...rule, pluginName: plugin.name})))
      .flatMap(rule =>
        rule.results.map(result => ({
          ...result,
          ruleName: rule.name,
          ruleErrorDescription: rule.errorDescription,
          pluginName: rule.pluginName,
        }))
      )
      .filter(result => !result.passed)
      .reduce<TAnalyzerErrorsBySpan>((prev, curr) => {
        const value = prev[curr.spanId] || [];
        return {...prev, [curr.spanId]: [...value, curr]};
      }, {});
  },

  getTestSpecsHashedBySpan(assertionResults: AssertionResults): TTestSpecsBySpan {
    return assertionResults.resultList
      .flatMap(assertionResult =>
        assertionResult.resultList.map(assertion => ({
          ...assertion,
          selector: assertionResult.selector || 'All Spans',
        }))
      )
      .flatMap(assertion =>
        assertion.spanResults.map(spanResult => ({
          ...spanResult,
          selector: assertion.selector,
          assertion: assertion.assertion,
        }))
      )
      .reduce<TTestSpecsBySpan>((prev, curr) => {
        const value = prev[curr.spanId] || {failed: [], passed: []};

        if (curr.passed) {
          value.passed.push(curr);
        } else {
          value.failed.push(curr);
        }

        return {...prev, [curr.spanId]: value};
      }, {});
  },

  getTestOutputsHashedBySpan(outputs: TestRunOutput[] = []): TTestOutputsBySpan {
    return outputs.reduce<TTestOutputsBySpan>((prev, curr) => {
      const value = prev[curr.spanId] || [];
      return {...prev, [curr.spanId]: [...value, curr]};
    }, {});
  },
});

export default TestRunService();
