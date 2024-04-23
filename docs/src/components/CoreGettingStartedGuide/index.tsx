/* eslint-disable global-require */

import React from 'react';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const CoreGettingStartedGuides = [
  {
    name: '👇 Install Tracetest Core',
    url: '/core/getting-started/installation',
    description: (
      <Translate >
        Set up Tracetest Core and start trace-based testing your distributed system.
      </Translate>
    ),
  },
  {
    name: '🙌 Open Tracetest Core',
    url: '/core/getting-started/open',
    description: (
      <Translate>
        After installing it, open Tracetest Core start to creating trace-based tests.
      </Translate>
    ),
  },
  {
    name: '🤔 Don\'t have OpenTelemetry?',
    url: '/getting-started/no-otel',
    description: (
      <Translate >
        Install OpenTelemetry in 5 minutes without any code changes!
      </Translate>
    ),
  },
  {
    name: '🤩 Open Source',
    url: 'https://github.com/kubeshop/tracetest',
    description: (
      <Translate>
        Check out the Tracetest GitHub repo! Please consider giving us a star! ⭐️
      </Translate>
    ),
  },
];

interface Props {
  name: string;
  url: string;
  description: JSX.Element;
}

function CoreGettingStartedGuideCard({name, url, description}: Props) {
  return (
    <div className="col col--6 margin-bottom--lg">
      <div className="gs__card">
      <div className="card">
        <Link to={url}>
          <div className="card__body">
            <Heading as="h3">{name}</Heading>
            <p>{description}</p>
          </div>
        </Link>
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
