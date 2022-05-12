import faker from '@faker-js/faker';
import {ComponentStory, ComponentMeta} from '@storybook/react';
import {TestState} from '../../../constants/TestRunResult.constants';

import ResultCard from '../ResultCard';

export default {
  title: 'Result Card',
  component: ResultCard,
  argTypes: {onClick: {action: 'clicked'}, onDelete: {action: 'deleted'}},
} as ComponentMeta<typeof ResultCard>;

const Template: ComponentStory<typeof ResultCard> = args => <ResultCard {...args} />;

export const Default = Template.bind({});
Default.args = {
  result: {
    resultId: faker.datatype.uuid(),
    testId: faker.datatype.uuid(),
    traceId: faker.datatype.uuid(),
    spanId: faker.datatype.uuid(),
    createdAt: faker.date.recent().toISOString(),
    completedAt: faker.date.recent().toISOString(),
    response: {},
    state: TestState.FINISHED,
    assertionResultState: false,
    assertionResult: [],
    totalAssertionCount: 10,
    passedAssertionCount: 7,
    failedAssertionCount: 3,
    executionTime: faker.datatype.number({min: 5, max: 25}),
  },
};
