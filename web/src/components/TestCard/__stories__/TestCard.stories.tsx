import faker from '@faker-js/faker';
import {ComponentStory, ComponentMeta} from '@storybook/react';
import {HTTP_METHOD} from '../../../constants/Common.constants';
import TestDefinition from '../../../models/TestDefinition.model';

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
    id: faker.datatype.uuid(),
    definition: TestDefinition({}),
    name: `${faker.name.firstName()} ${faker.name.lastName()}`,
    description: faker.lorem.sentences(),
    serviceUnderTest: {
      request: {
        url: faker.internet.url(),
        method: faker.internet.httpMethod().toUpperCase() as HTTP_METHOD,
      },
    },
  },
};
