/* eslint-disable global-require */

import React from 'react';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const TracetestCoreGettingStartedGuides = [
  {
    name: 'Option 3: I have a non-production environment and want to try trace-based testing',
    url: '/core/getting-started/overview',
    title: 'Hobby Self-hosted Tracetest Core (Open Source)',
    description: (
      <Translate>
        Deploy a hobby instance in your own infrastructure with Docker or Kubernetes. Not suitable for production workloads.
      </Translate>
    ),
  },
];

const TracetestCloudGettingStartedGuides = [
  {
    name: 'Option 1: I want to create tests without managing infrastructure (Most Popular)',
    url: '/getting-started/create-tracetest-account',
    title: 'Cloud-based Managed Tracetest (Free to get started!)',
    description: (
      <Translate>
        Use managed infrastructure with collaboration for teams, and additional features on top of Hobby Tracetest Core.
      </Translate>
    ),
  },
];

const TracetestOnPremGettingStartedGuides = [
  {
    name: 'Option 2: I need something secure, controlled, in my own infrastructure',
    url: 'https://tracetest.io/on-prem-installation',
    title: 'Enterprise Self-hosted Tracetest (Free trial, no credit card required)',
    description: (
      <Translate>
        Same experience as with Cloud-based Managed Tracetest but self-hosted in your own infrastructure.
      </Translate>
    ),
  },
];

interface Props {
  name: string;
  url: string;
  title: string;
  description: JSX.Element;
}

function GettingStartedGuideCard({name, url, title, description}: Props) {
  return (
    <div className="col col--12">
      <div className="gs__card">
      <div className="card">
        <Link to={url}>
          <div className="card__body">
            <Heading as="h3">{name}</Heading>
            <p>
              <b>{title}:&nbsp;</b>
              {description}
            </p>
          </div>
        </Link>
      </div>
      </div>
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

export function TracetestCloudGettingStartedGuidesRow(): JSX.Element {
  return (
    <div className="row">
      {TracetestCloudGettingStartedGuides.map((gettingStartedGuide) => (
        <GettingStartedGuideCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}

export function TracetestOnPremGettingStartedGuidesRow(): JSX.Element {
  return (
    <div className="row">
      {TracetestOnPremGettingStartedGuides.map((gettingStartedGuide) => (
        <GettingStartedGuideCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}
