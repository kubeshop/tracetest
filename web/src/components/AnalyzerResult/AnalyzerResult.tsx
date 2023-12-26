import BetaBadge from 'components/BetaBadge/BetaBadge';
import Link from 'components/Link';
import {COMMUNITY_SLACK_URL, OCTOLIINT_ISSUE_URL} from 'constants/Common.constants';
import LinterResult from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from './AnalyzerResult.styled';
import Empty from './Empty';
import GlobalResult from './GlobalResult';
import Plugins from './Plugins';

interface IProps {
  result: LinterResult;
  trace: Trace;
}

const AnalyzerResult = ({result: {score, minimumScore, plugins = [], passed}, trace}: IProps) => {
  const {linter} = useSettingsValues();

  return (
    <S.Container>
      <S.Title level={2}>
        Analyzer Results <BetaBadge />
      </S.Title>

      <S.Description>
        The Tracetest Analyzer is a plugin based framework used to analyze OpenTelemetry traces to help teams improve
        their instrumentation data, find potential problems and provide tips to fix the problems.{' '}
        {linter.enabled && (
          <>
            It can be globally disabled for all tests in <Link to="/settings?tab=analyzer">the settings page</Link>.{' '}
          </>
        )}
        We value your feedback on this beta release. Share your thoughts on <a href={COMMUNITY_SLACK_URL}>Slack</a> or add
        them to this <a href={OCTOLIINT_ISSUE_URL}>Issue</a>.
      </S.Description>
      {plugins.length ? (
        <>
          <GlobalResult score={score} minimumScore={minimumScore} allRulesPassed={passed} />
          <Plugins plugins={plugins} trace={trace} />
        </>
      ) : (
        <Empty />
      )}
    </S.Container>
  );
};

export default AnalyzerResult;
