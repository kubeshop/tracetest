import {ComponentStory, ComponentMeta} from '@storybook/react';
import TestMock from '../../../models/__mocks__/Test.mock';

import TestCard from '../TestCard';

export default {
  title: 'Test Card',
  component: TestCard,
  argTypes: {onClick: {action: 'clicked'}, onDelete: {action: 'deleted'}, onRunTest: {action: 'runTest'}},
} as ComponentMeta<typeof TestCard>;

const Template: ComponentStory<typeof TestCard> = args => <TestCard {...args} />;

export const Default = Template.bind({});
Default.args = {
  test: TestMock.model(),
};
