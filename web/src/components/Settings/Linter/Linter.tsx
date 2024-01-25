import {ANALYZER_DOCUMENTATION_URL, COMMUNITY_SLACK_URL, OCTOLIINT_ISSUE_URL} from 'constants/Common.constants';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import LinterForm from './LinterForm';
import * as S from '../common/Settings.styled';

const Linter = () => (
  <S.Container>
    <S.Title level={2}>Analyzer (Beta)</S.Title>
    <S.Description>
      <p>
        The Tracetest Analyzer is a plugin based framework used to audit OpenTelemetry traces to help teams improve
        their instrumentation data, find potential problems and provide tips to fix the problems. We have released this
        initial version to get feedback from the community. Have thoughts about how to improve the Tracetest Analyzer?
        Add to this <a href={OCTOLIINT_ISSUE_URL}>Issue</a> or call us on <a href={COMMUNITY_SLACK_URL}>Slack</a>!
      </p>
      <DocsBanner>
        Need more information about the Analyzer?{' '}
        <a href={ANALYZER_DOCUMENTATION_URL} target="_blank">
          Learn more in our docs
        </a>
      </DocsBanner>
    </S.Description>
    <S.FormContainer>
      <LinterForm />
    </S.FormContainer>
  </S.Container>
);

export default Linter;
