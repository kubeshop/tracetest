import {TRawLinterResultPluginRuleResultError} from 'models/LinterResult.model';
import TestRunOutput from 'models/TestRunOutput.model';
import {TTestSchemas} from './Common.types';

export enum RunErrorTypes {
  MissingVariables = 'missingVariables',
  Unknown = 'unknown',
}

export type TTestRunState = NonNullable<TTestSchemas['TestRun']['state'] | 'WAITING' | 'SKIPPED' | 'FAILED'>;

/* AnalyzerErrors by Span types */

export type TAnalyzerError = {
  ruleName: string;
  ruleErrorDescription: string;
  pluginName: string;
  passed: boolean;
  spanId: string;
  errors: TRawLinterResultPluginRuleResultError[];
  severity: 'error' | 'warning';
};

export type TAnalyzerErrorsBySpan = Record<string, TAnalyzerError[]>;

/*  TestSpecs by Span types */

export type TTestSpec = {
  selector: string;
  assertion: string;
  spanId: string;
  observedValue: string;
  passed: boolean;
  error: string;
};

export type TTestSpecSummary = {
  failed: TTestSpec[];
  passed: TTestSpec[];
};

export type TTestSpecsBySpan = Record<string, TTestSpecSummary>;

/*  TestOutputs by Span types */

export type TTestOutputsBySpan = Record<string, TestRunOutput[]>;
