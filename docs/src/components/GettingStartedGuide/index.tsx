/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Translate from '@docusaurus/Translate';
import Link from '@docusaurus/Link';
import Heading from '@theme/Heading';

const GettingStartedGuides = [
  {
    name: '👇 Install Tracetest',
    url: './getting-started/installation',
    description: (
      <Translate >
        Set up Tracetest and start trace-based testing your distributed system.
      </Translate>
    ),
    button: 'Set up Tracetest',
  },
  {
    name: '🙌 Open Tracetest',
    url: './getting-started/open',
    description: (
      <Translate>
        After installing it, open Tracetest start to creating trace-based tests.
      </Translate>
    ),
    button: 'Create tests',
  },
  {
    name: '🤔 Don\'t have OpenTelemetry?',
    url: './getting-started/no-otel',
    description: (
      <Translate >
        Install OpenTelemetry in 5 minutes without any code changes!
      </Translate>
    ),
    button: 'Set up OTel',
  },
  {
    name: '🤩 Open Source',
    url: 'https://github.com/kubeshop/tracetest',
    description: (
      <Translate>
        Check out the Tracetest GitHub repo! Please consider giving us a star! ⭐️
      </Translate>
    ),
    button: 'Go to GitHub',
  },
  {
    name: '⚙️ Configure trace data stores',
    url: '../configuration/overview#supported-trace-data-stores',
    description: (
      <Translate>
        Connect your existing trace data store or send traces to Tracetest directly!
      </Translate>
    ),
    button: 'Configure Tracetest',
  },
  {
    name: '🙄 New to Trace-based Testing?',
    url: '../concepts/what-is-trace-based-testing',
    description: (
      <Translate>
        Read about the concepts of trace-based testing to learn more!
      </Translate>
    ),
    button: 'View Concepts',
  },
];

interface Props {
  name: string;
  url: string;
  button: string;
  description: JSX.Element;
}

function GettingStartedGuideCard({name, url, description, button}: Props) {
  return (
    <div className="col col--6 margin-bottom--lg">
      <div className={clsx('card')}>
        <div className="card__body">
          <Heading as="h3">{name}</Heading>
          <p>{description}</p>
        </div>
        <div className="card__footer">
          <div className="button-group button-group--block">
            <Link className="button button--secondary" to={url}>
              {button}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

export function GettingStartedGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {GettingStartedGuides.map((gettingStartedGuide) => (
        <GettingStartedGuideCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}
