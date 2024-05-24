/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const ExamplesTutorialsOverview = [
  {
    name: '🍱 Recipes',
    url: '/examples-tutorials/recipes',
    description: (
      <Translate>
        Self-contained guides to popular use cases.
      </Translate>
    ),
    button: 'Start building',
  },
  {
    name: '🛠️ Tools & Integrations',
    url: '/tools-and-integrations/overview',
    description: (
      <Translate>
        Tools and integrations examples with Tracetest.
      </Translate>
    ),
    button: 'Integrate',
  },
  {
    name: '⚙️ CI/CD Automation',
    url: '/ci-cd-automation/overview',
    description: (
      <Translate >
        Run Tracetest in a CI/CD pipeline!
      </Translate>
    ),
    button: 'Automate',
  },
  {
    name: '📽️ Videos & Recordings',
    url: '/examples-tutorials/videos',
    description: (
      <Translate>
        Tracetest video guides and conference talks!
      </Translate>
    ),
    button: 'Watch now',
  },
  {
    name: '🎙️ Webinars',
    url: '/examples-tutorials/webinars',
    description: (
      <Translate >
        Watch on-demand live streams and community calls!
      </Translate>
    ),
    button: 'Watch now',
  },
  {
    name: '👨‍💻 Blog Post Tutorials',
    url: '/examples-tutorials/tutorials',
    description: (
      <Translate >
        Check out blog posts with Tracetest-related content.
      </Translate>
    ),
    button: 'Learn more',
  },
];

interface Props {
  name: string;
  url: string;
  button: string;
  description: JSX.Element;
}

function ExamplesTutorialsOverviewCard({name, url, description, button}: Props) {
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

export function ExamplesTutorialsOverviewRow(): JSX.Element {
  return (
    <div className="row">
      {ExamplesTutorialsOverview.map((gettingStartedGuide) => (
        <ExamplesTutorialsOverviewCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}
