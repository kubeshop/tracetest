import {DISCORD_URL, OCTOLIINT_ISSUE_URL} from 'constants/Common.constants';
import LinterForm from './LinterForm';
import * as S from '../common/Settings.styled';

const Linter = () => (
  <S.Container>
    <S.Description>
      <p>
        The Tracetest Analyzer is a plugin based framework used to audit OpenTelemetry traces to help teams improve
        their instrumentation data, find potential problems and provide tips to fix the problems. We have released this
        initial version to get feedback from the community. Have thoughts about how to improve the Tracetest Analyzer?
        Add to this <a href={OCTOLIINT_ISSUE_URL}>Issue</a> or <a href={DISCORD_URL}>Discord</a>! Currently, the
        analyzer supports the following plugins:
      </p>
      <S.LinterPluginList>
        <li>
          <b>OpenTelemetry Semantic Conventions.</b> Enforce standards for spans and attributes
        </li>
        <li>
          <b>Security.</b> Enforce security for spans and attributes
        </li>
        <li>
          <b>Common problems.</b> Helps you find common problems with your application
        </li>
      </S.LinterPluginList>
    </S.Description>
    <S.FormContainer>
      <LinterForm />
    </S.FormContainer>
  </S.Container>
);

export default Linter;
