import faker from '@faker-js/faker';
import {ComponentStory, ComponentMeta} from '@storybook/react';
import {HTTP_METHOD} from '../../../constants/Common.constants';

import TestCard from '../TestCard';

export default {
  title: 'Test Card',
  component: TestCard,
  argTypes: {onClick: {action: 'clicked'}, onDelete: {action: 'deleted'}, onRunTest: {action: 'runTest'}},
} as ComponentMeta<typeof TestCard>;

const Template: ComponentStory<typeof TestCard> = args => <TestCard {...args} />;

export const Default = Template.bind({});
Default.args = {
  test: {
    testId: faker.datatype.uuid(),
    name: `${faker.name.firstName()} ${faker.name.lastName()}`,
    description: faker.lorem.sentences(),
    serviceUnderTest: {
      id: faker.datatype.uuid(),
      request: {
        url: faker.internet.url(),
        method: faker.internet.httpMethod().toUpperCase() as HTTP_METHOD,
      },
    },
    assertions: [],
  },
};
