import {PostHogProvider} from 'posthog-js/react';
import Env from 'utils/Env';
import Content from './Content';

const posthogKey = Env.get('posthogKey');

const options = {
  api_host: 'https://app.posthog.com',
  debug: true,
  opt_out_capturing_by_default: true,
};

interface IProps {
  children: React.ReactNode;
}

const CaptureWrapper = ({children}: IProps) => (
  <PostHogProvider apiKey={posthogKey} options={options}>
    <Content>{children}</Content>
  </PostHogProvider>
);

export default CaptureWrapper;
