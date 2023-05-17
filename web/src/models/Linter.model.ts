import {Model, TLintersSchemas} from 'types/Common.types';

export type TRawLinter = TLintersSchemas['LinternResource'];
type Linter = Model<Model<TRawLinter, {}>['spec'], {}>;

type TRawLinterPlugin = TLintersSchemas['LinternResourcePlugin'];
type LinterPlugin = Model<TRawLinterPlugin, {}>;

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
