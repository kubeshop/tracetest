/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const CoreGettingStartedGuides = [
  {
    name: '👇 Install Tracetest Core',
    url: './installation',
    description: (
      <Translate >
        Set up Tracetest Core and start trace-based testing your distributed system.
      </Translate>
    ),
    button: 'Start',
  },
  {
    name: '🙌 Open Tracetest Core',
    url: './open',
    description: (
      <Translate>
        After installing it, open Tracetest Core start to creating trace-based tests.
      </Translate>
    ),
    button: 'Create tests',
  },
  {
    name: '🤔 Don\'t have OpenTelemetry?',
    url: '/getting-started/no-otel',
    description: (
      <Translate >
        Install OpenTelemetry in 5 minutes without any code changes!
      </Translate>
    ),
    button: 'Configure OpenTelemetry',
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
];

interface Props {
  name: string;
  url: string;
  button: string;
  description: JSX.Element;
}

function CoreGettingStartedGuideCard({name, url, description, button}: Props) {
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

export function CoreGettingStartedGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {CoreGettingStartedGuides.map((coreGettingStartedGuide) => (
        <CoreGettingStartedGuideCard key={coreGettingStartedGuide.name} {...coreGettingStartedGuide} />
      ))}
    </div>
  );
}
