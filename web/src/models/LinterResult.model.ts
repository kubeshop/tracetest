import {Model, TLintersSchemas} from 'types/Common.types';

type TRawLinterResult = TLintersSchemas['LinterResult'];
type TRawLinterResultPlugin = TLintersSchemas['LinterResultPlugin'];
type TRawLinterResultPluginRule = TLintersSchemas['LinterResultPluginRule'];
type TRawLinterResultPluginRuleResult = TLintersSchemas['LinterResultPluginRuleResult'];
type TRawLinterResultPluginRuleResultGroupedError = TLintersSchemas['LinterResultPluginRuleResultGroupedError'];

type LinterResultPluginRuleResultGroupedError = Model<TRawLinterResultPluginRuleResultGroupedError, {}>;
type LinterResultPluginRuleResult = Model<TRawLinterResultPluginRuleResult, {}>;
type LinterResultPluginRule = Model<TRawLinterResultPluginRule, {results: LinterResultPluginRuleResult[]}>;
type LinterResultPlugin = Model<TRawLinterResultPlugin, {rules: LinterResultPluginRule[]}>;
type LinterResult = Model<TRawLinterResult, {plugins: LinterResultPlugin[]; isFailed: boolean}>;

function LinterResultPluginRuleResultGroupedError({
  error = '',
  values = [],
}: TRawLinterResultPluginRuleResultGroupedError = {}): LinterResultPluginRuleResultGroupedError {
  return {error, values};
}

function LinterResultPluginRuleResult({
  spanId = '',
  errors = [],
  groupedErrors = [],
  passed = false,
  severity = 'error',
}: TRawLinterResultPluginRuleResult = {}): LinterResultPluginRuleResult {
  return {
    spanId,
    errors,
    groupedErrors: groupedErrors.map(groupedError => LinterResultPluginRuleResultGroupedError(groupedError)),
    passed,
    severity,
  };
}

function LinterResultPluginRule({
  name = '',
  description = '',
  passed = false,
  weight = 0,
  tips = [],
  results = [],
}: TRawLinterResultPluginRule = {}): LinterResultPluginRule {
  return {
    name,
    description,
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
