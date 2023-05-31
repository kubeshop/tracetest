import {DISCORD_URL, OCTOLIINT_ISSUE_URL} from 'constants/Common.constants';
import LinterForm from './LinterForm';
import * as S from '../common/Settings.styled';

const Linter = () => (
  <S.Container>
    <S.Description>
      <p>
        This beta release of the Tracetest Linter its a plugin based framework used to analyze Open Telemetry traces to
        help teams improve their instrumentation data, find potential problems and provide tips to fix the problems. We
        have released this initial version to get feedback from the community. Have thoughts about how to improve the
        Tracetest Linter? Add to this
        <a href={OCTOLIINT_ISSUE_URL}>Issue</a> or <a href={DISCORD_URL}>Discord</a>! Currently, the linter supports the
        following plugins:
      </p>
      <S.LinterPluginList>
        <li>
          <b>Open Telemetry Semantic Conventions.</b> Enforce standards for spans and attributes
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
