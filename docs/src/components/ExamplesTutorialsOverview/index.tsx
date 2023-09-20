/* eslint-disable global-require */

import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import Translate from '@docusaurus/Translate';
import Heading from '@theme/Heading';

const ExamplesTutorialsOverview = [
  {
    name: '👨‍💻 Tutorials',
    url: '/examples-tutorials/tutorials',
    description: (
      <Translate >
        Check out the following blog posts with Tracetest-related content.
      </Translate>
    ),
    button: 'Learn more',
  },
  {
    name: '🍱 Recipes',
    url: '/examples-tutorials/recipes',
    description: (
      <Translate>
        Short, self-contained, and runnable solutions to popular use cases.
      </Translate>
    ),
    button: 'Start building',
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
    name: '🛠️ Tools & Integrations',
    url: '/tools-and-integrations/overview',
    description: (
      <Translate>
        Check out tools and integrations with Tracetest.
      </Translate>
    ),
    button: 'Integrate',
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
    name: '📽️ Videos',
    url: '/examples-tutorials/videos',
    description: (
      <Translate>
        Check out Tracetest video guides and conference talks!
      </Translate>
    ),
    button: 'Watch now',
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

export function ExamplesTutorialsOverviewRow(): JSX.Element {
  return (
    <div className="row">
      {ExamplesTutorialsOverview.map((gettingStartedGuide) => (
        <ExamplesTutorialsOverviewCard key={gettingStartedGuide.name} {...gettingStartedGuide} />
      ))}
    </div>
  );
}
