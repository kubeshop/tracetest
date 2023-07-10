import {Model, TLintersSchemas} from 'types/Common.types';

export type TRawLinter = TLintersSchemas['LinterResource'];
type Linter = Model<Model<TRawLinter, {}>['spec'], {plugins: LinterPlugin[]}>;

type TRawLinterPlugin = TLintersSchemas['LinterResourcePlugin'];
export type LinterPlugin = Model<
  TRawLinterPlugin,
  {
    rules: LinterRule[];
  }
>;

type TRawLinterRule = TLintersSchemas['LinterResourceRule'];
export type LinterRule = Model<TRawLinterRule, {}>;

export function LinterPlugin({
  name = '',
  id = '',
  enabled = false,
  rules = [],
  description = '',
}: TRawLinterPlugin = {}): LinterPlugin {
  return {name, id, enabled, description, rules: rules.map(rule => LinterRule(rule))};
}

export enum LinterRuleErrorLevel {
  ERROR = 'error',
  WARNING = 'warning',
  DISABLED = 'disabled',
}

export function LinterRule({
  id = '',
  weight = 0,
  errorLevel = 'error',
  description = '',
  errorDescription = '',
  tips = [],
  documentation = '',
  name = '',
}: TRawLinterRule = {}): LinterRule {
  return {id, weight, errorLevel, name, description, errorDescription, tips, documentation};
}

function Linter({
  spec: {id = '', name = '', enabled = false, minimumScore = 100, plugins = []} = {},
}: TRawLinter = {}): Linter {
  return {
    id,
    name,
    enabled,
    minimumScore,
    plugins: plugins.map(plugin => LinterPlugin(plugin)),
  };
}

export default Linter;
