import {Link} from 'react-router-dom';
import {Col, Row, Tooltip} from 'antd';
import LinterResult from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import {DISCORD_URL, OCTOLIINT_ISSUE_URL} from 'constants/Common.constants';
import {useSettingsValues} from 'providers/SettingsValues/SettingsValues.provider';
import * as S from './AnalyzerResult.styled';
import LintScore from '../LintScore/LintScore';
import BetaBadge from '../BetaBadge/BetaBadge';
import Empty from './Empty';
import Plugins from './Plugins';

interface IProps {
  result: LinterResult;
  trace: Trace;
}

const AnalyzerResult = ({result: {score, passed, plugins = []}, trace}: IProps) => {
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
          <Row gutter={[16, 16]}>
            <Col span={8} key="avg_result">
              <Tooltip title="Tracetest core system supports analyzer evaluation as part of the testing capabilities.">
                <S.ScoreContainer>
                  <S.Subtitle level={3}>Trace Analyzer Result</S.Subtitle> <LintScore score={score} passed={passed} />
                </S.ScoreContainer>
              </Tooltip>
            </Col>
            {plugins.map(plugin => (
              <Col span={8} key={plugin.name}>
                <Tooltip title={plugin.description}>
                  <S.ScoreContainer key={plugin.name}>
                    <S.Subtitle level={3}>{plugin.name}</S.Subtitle>
                    <LintScore score={plugin.score} passed={plugin.passed} />
                  </S.ScoreContainer>
                </Tooltip>
              </Col>
            ))}
          </Row>
          <Plugins plugins={plugins} trace={trace} />
        </>
      ) : (
        <Empty />
      )}
    </S.Container>
  );
};

export default AnalyzerResult;
