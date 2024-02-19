/* eslint-disable global-require */

import React from 'react';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const TracetestGettingStartedGuides = [
  {
    name: 'Tracetest 🚀',
    url: '/getting-started/installation',
    description: (
      <Translate >
        Set up Tracetest and start trace-based testing your distributed system.
      </Translate>
    ),
    button: 'Start',
  },
];

const TracetestCoreGettingStartedGuides = [
  {
    name: 'Tracetest Core 🪨 ',
    url: '/core/getting-started/installation',
    description: (
      <Translate>
        Use the open-source Tracetest Core in your own infrastructure.
      </Translate>
    ),
    button: 'Go to Core',
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

export function TracetestGettingStartedGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {TracetestGettingStartedGuides.map((gettingStartedGuide) => (
        <GettingStartedGuideCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}

export function TracetestCoreGettingStartedGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {TracetestCoreGettingStartedGuides.map((gettingStartedGuide) => (
        <GettingStartedGuideCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}
