import faker from '@faker-js/faker';
import {ComponentStory, ComponentMeta} from '@storybook/react';
import {HTTP_METHOD} from '../../../constants/Common.constants';
import { TestState } from '../../../constants/TestRun.constants';
import TestDefinition from '../../../models/TestDefinition.model';

import TestHeader from '../TestHeader';

export default {
  title: 'Test Header',
  component: TestHeader,
  argTypes: {onBack: {action: 'onBack'}},
} as ComponentMeta<typeof TestHeader>;

const Template: ComponentStory<typeof TestHeader> = args => <TestHeader {...args} />;

export const Default = Template.bind({});
Default.args = {
  test: {
    id: faker.datatype.uuid(),
    name: `${faker.name.firstName()} ${faker.name.lastName()}`,
    description: faker.lorem.sentences(),
    definition: TestDefinition({}),
    serviceUnderTest: {
      request: {
        url: faker.internet.url(),
        method: faker.internet.httpMethod().toUpperCase() as HTTP_METHOD,
      },
    },
  },
};

export const WithState = Template.bind({});
WithState.args = {
  test: {
    id: faker.datatype.uuid(),
    name: `${faker.name.firstName()} ${faker.name.lastName()}`,
    description: faker.lorem.sentences(),
    definition: TestDefinition({}),
    serviceUnderTest: {
      request: {
        url: faker.internet.url(),
        method: faker.internet.httpMethod().toUpperCase() as HTTP_METHOD,
      },
    },
  },
  testState: TestState.FINISHED,
};
