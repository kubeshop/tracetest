import {Model, TLintersSchemas} from 'types/Common.types';

type TRawLinterResult = TLintersSchemas['LinternResult'];
type TRawLinterResultPlugin = TLintersSchemas['LinternResultPlugin'];
type TRawLinterResultPluginRule = TLintersSchemas['LinternResultPluginRule'];
type TRawLinterResultPluginRuleResult = TLintersSchemas['LinternResultPluginRuleResult'];

type LinterResultPluginRuleResult = Model<TRawLinterResultPluginRuleResult, {}>;
type LinterResultPluginRule = Model<TRawLinterResultPluginRule, {results: LinterResultPluginRuleResult[]}>;
type LinterResultPlugin = Model<TRawLinterResultPlugin, {rules: LinterResultPluginRule[]}>;
type LinterResult = Model<TRawLinterResult, {plugins: LinterResultPlugin[]}>;

function LinterResultPluginRuleResult({
  spanId = '',
  errors = [],
  passed = false,
  severity = 'warning',
}: TRawLinterResultPluginRuleResult = {}): LinterResultPluginRuleResult {
  return {spanId, errors, passed, severity};
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

function LinterResult({passed = false, score = 0, plugins = []}: TRawLinterResult = {}): LinterResult {
  return {
    passed,
    score,
    plugins: plugins.map(plugin => LinterResultPlugin(plugin)),
  };
}

export default LinterResult;
