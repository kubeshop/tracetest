/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Translate from '@docusaurus/Translate';
import Link from '@docusaurus/Link';
import Heading from '@theme/Heading';

const WelcomeGuides = [
  {
    name: '👇 Getting Started',
    url: './getting-started/overview',
    description: (
      <Translate >
        Set up Tracetest and start trace-based testing your distributed system.
      </Translate>
    ),
    button: 'Start',
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
    url: '/configuration/overview#supported-trace-data-stores',
    description: (
      <Translate>
        Connect your existing trace data store or send traces to Tracetest directly!
      </Translate>
    ),
    button: 'Configure',
  },
  {
    name: '🙄 New to Trace-based Testing?',
    url: '/concepts/what-is-trace-based-testing',
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

function WelcomeGuideCard({name, url, description, button}: Props) {
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

export function WelcomeGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {WelcomeGuides.map((WelcomeGuide) => (
        <WelcomeGuideCard key={WelcomeGuide.name} {...WelcomeGuide} />
      ))}
    </div>
  );
}
