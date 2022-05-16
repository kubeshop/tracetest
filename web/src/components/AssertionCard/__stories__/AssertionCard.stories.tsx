import faker from '@faker-js/faker';
import {ComponentStory, ComponentMeta} from '@storybook/react';
import {LOCATION_NAME} from '../../../constants/Span.constants';
import {SELECTOR_DEFAULT_ATTRIBUTES} from '../../../constants/SemanticGroupNames.constants';

import AssertionCard from '../AssertionCard';
import SpanMock from '../../../models/__mocks__/Span.mock';

export default {
  title: 'Assertion Card',
  component: AssertionCard,
  argTypes: {onSelectSpan: {action: 'noSelectSpan'}},
} as ComponentMeta<typeof AssertionCard>;

const Template: ComponentStory<typeof AssertionCard> = args => <AssertionCard {...args} />;

export const Default = Template.bind({});
Default.args = {
  assertionResult: {
    assertion: {
      assertionId: faker.datatype.uuid(),
      selectors: [
        {
          locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
          propertyName: faker.random.arrayElement(SELECTOR_DEFAULT_ATTRIBUTES[0].attributes),
          value: 'http',
          valueType: faker.random.word(),
        },
        {
          locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
          propertyName: faker.random.arrayElement(SELECTOR_DEFAULT_ATTRIBUTES[0].attributes),
          value: 'pokeshop',
          valueType: faker.random.word(),
        },
      ],
      spanAssertions: [],
    },
    spanListAssertionResult: faker.datatype.array(faker.datatype.number({min: 1, max: 10})).map(() => ({
      span: SpanMock.model(),
      resultList: [],
    })),
  },
};
