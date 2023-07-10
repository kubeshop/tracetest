import {Model, TLintersSchemas} from 'types/Common.types';
import {LinterRuleErrorLevel} from './Linter.model';

type TRawLinterResult = TLintersSchemas['LinterResult'];
type TRawLinterResultPlugin = TLintersSchemas['LinterResultPlugin'];
type TRawLinterResultPluginRule = TLintersSchemas['LinterResultPluginRule'];
type TRawLinterResultPluginRuleResult = TLintersSchemas['LinterResultPluginRuleResult'];
export type TRawLinterResultPluginRuleResultError = TLintersSchemas['LinterResultPluginRuleResultError'];

type LinterResultPluginRuleResultError = Model<TRawLinterResultPluginRuleResultError, {}>;
type LinterResultPluginRuleResult = Model<TRawLinterResultPluginRuleResult, {}>;
export type LinterResultPluginRule = Model<
  TRawLinterResultPluginRule,
  {results: LinterResultPluginRuleResult[]; level: LinterRuleErrorLevel}
>;
export type LinterResultPlugin = Model<TRawLinterResultPlugin, {rules: LinterResultPluginRule[]}>;
type LinterResult = Model<TRawLinterResult, {plugins: LinterResultPlugin[]; isFailed: boolean}>;

function LinterResultPluginRuleResultError({
  value = '',
  expected = '',
  description = '',
  suggestions = [],
}: TRawLinterResultPluginRuleResultError = {}): LinterResultPluginRuleResultError {
  return {value, expected, description, suggestions};
}

function LinterResultPluginRuleResult({
  spanId = '',
  errors = [],
  passed = false,
  severity = 'error',
}: TRawLinterResultPluginRuleResult = {}): LinterResultPluginRuleResult {
  return {
    spanId,
    errors: errors.map(error => LinterResultPluginRuleResultError(error)),
    passed,
    severity,
  };
}

function LinterResultPluginRule({
  name = '',
  description = '',
  errorDescription = '',
  passed = false,
  weight = 0,
  tips = [],
  results = [],
  level = 'error',
}: TRawLinterResultPluginRule = {}): LinterResultPluginRule {
  return {
    name,
    level: level as LinterRuleErrorLevel,
    description,
    errorDescription,
    passed,
    weight,
    tips,
    results: results.map(result => LinterResultPluginRuleResult(result)),
  };
}

function LinterResultPlugin({
  name = '',
  description = '',
  passed = false,
  score = 0,
  rules = [],
}: TRawLinterResultPlugin = {}): LinterResultPlugin {
  return {name, description, passed, score, rules: rules.map(rule => LinterResultPluginRule(rule))};
}

function LinterResult({
  passed = false,
  score = 0,
  plugins = [],
  minimumScore = 0,
}: TRawLinterResult = {}): LinterResult {
  return {
    passed,
    score,
    minimumScore,
    plugins: plugins.map(plugin => LinterResultPlugin(plugin)),
    isFailed: score < minimumScore,
  };
}

export default LinterResult;
