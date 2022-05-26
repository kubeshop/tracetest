import {ComponentStory, ComponentMeta} from '@storybook/react';
import {TestState} from '../../../constants/TestRun.constants';
import TestMock from '../../../models/__mocks__/Test.mock';

import TestHeader from '../TestHeader';

export default {
  title: 'Test Header',
  component: TestHeader,
  argTypes: {onBack: {action: 'onBack'}},
} as ComponentMeta<typeof TestHeader>;

const Template: ComponentStory<typeof TestHeader> = args => <TestHeader {...args} />;

export const Default = Template.bind({});
Default.args = {
  test: TestMock.model(),
};

export const WithState = Template.bind({});
WithState.args = {
  test: TestMock.model(),
  testState: TestState.FINISHED,
};
