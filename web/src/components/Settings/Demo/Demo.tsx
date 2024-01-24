import DocsBanner from 'components/DocsBanner';
import {DEMO_DOCUMENTATION_URL, OTEL_DEMO_GITHUB, POKESHOP_GITHUB} from 'constants/Common.constants';
import DemoForm from './DemoForm';
import * as S from '../common/Settings.styled';

const Demo = () => (
  <S.Container>
    <S.Title level={2}>Demo</S.Title>
    <S.Description>
      <p>
        Tracetest has the option to enable Test examples for our{' '}
        <a href={POKESHOP_GITHUB} target="_blank">
          Pokeshop Demo App
        </a>{' '}
        or the{' '}
        <a href={OTEL_DEMO_GITHUB} target="_blank">
          OpenTelemetry Astronomy Shop Demo
        </a>
        . You will require an instance of those applications running alongside your Tracetest server to be able to use
        them.
      </p>
      <DocsBanner>
        Need more information about the Demos?{' '}
        <a href={DEMO_DOCUMENTATION_URL} target="__blank">
          Learn more in our docs
        </a>
      </DocsBanner>
    </S.Description>
    <S.FormContainer>
      <DemoForm />
    </S.FormContainer>
  </S.Container>
);

export default Demo;
