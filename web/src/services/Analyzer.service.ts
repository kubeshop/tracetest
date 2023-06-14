import LinterResult from 'models/LinterResult.model';

const MAX_PLUGIN_SCORE = 100;

const AnalyzerService = () => ({
  getPlugins(plugins: LinterResult['plugins'], showOnlyErrors: boolean): LinterResult['plugins'] {
    return plugins
      .filter(plugin => !showOnlyErrors || plugin.score < MAX_PLUGIN_SCORE)
      .map(plugin => ({
        ...plugin,
        rules: plugin.rules
          .filter(rule => !showOnlyErrors || !rule.passed)
          .map(rule => ({...rule, results: rule?.results?.filter(result => !showOnlyErrors || !result.passed)})),
      }));
  },
});

export default AnalyzerService();
