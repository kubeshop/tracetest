import {Model, TLintersSchemas} from 'types/Common.types';

export type TRawLinter = TLintersSchemas['LinterResource'];
type Linter = Model<Model<TRawLinter, {}>['spec'], {plugins: LinterPlugin[]}>;

type TRawLinterPlugin = TLintersSchemas['LinterResourcePlugin'];
export type LinterPlugin = Model<TRawLinterPlugin, {}>;

function LinterPlugin({name = '', enabled = false, required = false}: TRawLinterPlugin = {}): LinterPlugin {
  return {name, enabled, required};
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
