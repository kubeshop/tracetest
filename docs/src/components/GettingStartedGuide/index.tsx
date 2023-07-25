/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Heading from '@theme/Heading';

const GettingStartedGuides = [
  {
    name: '👇 Install Tracetest',
    url: './installation',
    button: 'Set up Tracetest',
  },
  {
    name: '🙌 Open Tracetest',
    url: './open',
    button: 'Create tests',
  },
  {
    name: '🤔 Don\'t have OpenTelemetry?',
    url: './no-otel',
    button: 'Set up tracing',
  },
];

interface Props {
  name: string;
  url: string;
  button: string;
}

function GettingStartedGuideCard({name, url, button}: Props) {
  return (
    <div className="col col--4 margin-bottom--lg">
      <div className={clsx('card')}>
        <div className="card__body">
          <Heading as="h4" style={{ margin: 0 }}>{name}</Heading>
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
