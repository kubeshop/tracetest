import React from 'react';
import ReactDOM from 'react-dom/client';
import * as Sentry from '@sentry/react';
import {BrowserTracing} from '@sentry/tracing';
import './antd-theme/antd-customized.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import {IEnv} from './types/Common.types';
import {SENTRY_ALLOWED_URLS, SENTRY_DNS} from './constants/Common.constants';

declare global {
  interface Window {
    ENV: IEnv;
  }
}

Sentry.init({
  dsn: SENTRY_DNS,
  allowUrls: SENTRY_ALLOWED_URLS,
  integrations: [new BrowserTracing()],
  tracesSampleRate: 1.0,
});

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
