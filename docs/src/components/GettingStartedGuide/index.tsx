/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const GettingStartedGuides = [
  {
    name: '🌥️ Tracetest Cloud',
    url: './installation',
    description: (
      <Translate >
        Set up Tracetest and start trace-based testing your distributed system.
      </Translate>
    ),
    button: 'Start',
  },
  {
    name: '🪨 Tracetest Core',
    url: '../core/getting-started/installation',
    description: (
      <Translate>
        Use the open-source Tractest Core in your infrastructure.
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

export function GettingStartedGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {GettingStartedGuides.map((gettingStartedGuide) => (
        <GettingStartedGuideCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}
