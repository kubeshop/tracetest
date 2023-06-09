import {Link} from 'react-router-dom';
import LinterResult from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import {DISCORD_URL, OCTOLIINT_ISSUE_URL} from 'constants/Common.constants';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from './AnalyzerResult.styled';
import BetaBadge from '../BetaBadge/BetaBadge';
import Empty from './Empty';
import Plugins from './Plugins';
import GlobalResult from './GlobalResult';

interface IProps {
  result: LinterResult;
  trace: Trace;
}

const AnalyzerResult = ({result: {score, minimumScore, plugins = []}, trace}: IProps) => {
  const {linter} = useSettingsValues();

  return (
    <S.Container>
      <S.Title level={2}>
        Analyzer Results <BetaBadge />
      </S.Title>
      <S.Description>
        The Tracetest Analyzer is a plugin based framework used to analyze OpenTelemetry traces to help teams improve
        their instrumentation data, find potential problems and provide tips to fix the problems.{' '}
        {linter.enabled ? (
          <>
            If you want to disable the analyzer for all tests, go to the{' '}
            <Link to="/settings?tab=analyzer">settings page</Link>.
          </>
        ) : (
          ''
        )}
        We have released this initial version to get feedback from the community. Have thoughts about how to improve the
        Tracetest Analyzer? Add to this <a href={OCTOLIINT_ISSUE_URL}>Issue</a> or <a href={DISCORD_URL}>Discord</a>!
      </S.Description>

      {plugins.length ? (
        <>
          <GlobalResult score={score} minimumScore={minimumScore} />
          <Plugins plugins={plugins} trace={trace} />
        </>
      ) : (
        <Empty />
      )}
    </S.Container>
  );
};

export default AnalyzerResult;
