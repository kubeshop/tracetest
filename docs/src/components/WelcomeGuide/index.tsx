/* eslint-disable global-require */

import React from 'react';
import Translate from '@docusaurus/Translate';
import Link from '@docusaurus/Link';
import Heading from '@theme/Heading';

const WelcomeGuides = [
  {
    name: '👉 Getting Started',
    url: './getting-started/overview',
    description: (
      <Translate>
        Check out guides and installation examples to get started quickly!
      </Translate>
    ),
  },
  {
    name: '👇 Examples & Guides',
    url: './examples-tutorials/recipes',
    description: (
      <Translate >
        Hit the ground running with these trace-based testing guides!
      </Translate>
    ),
  },
  {
    name: '⚙️ Configure Access & Trace Ingestion',
    url: '/configuration/overview',
    description: (
      <Translate>
        Configure app access & connect tracing backend or OTLP ingestion!
      </Translate>
    ),
  },
  {
    name: '🙄 New to Trace-based Testing?',
    url: '/concepts/what-is-trace-based-testing',
    description: (
      <Translate>
        Read about the concepts of trace-based testing to learn more!
      </Translate>
    ),
  },
];

interface Props {
  name: string;
  url: string;
  description: JSX.Element;
}

function WelcomeGuideCard({name, url, description}: Props) {
  return (
    <div className="col col--6 margin-bottom--lg">
      <div className="w__card">
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

export function WelcomeGuideCardsRow(): JSX.Element {
  return (
    <div className="row">
      {WelcomeGuides.map((WelcomeGuide) => (
        <WelcomeGuideCard key={WelcomeGuide.name} {...WelcomeGuide} />
      ))}
    </div>
  );
}
