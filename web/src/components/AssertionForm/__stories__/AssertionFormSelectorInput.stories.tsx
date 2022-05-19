import {ComponentStory, ComponentMeta} from '@storybook/react';
import SpanMock from '../../../models/__mocks__/Span.mock';

import AssertionFormSelectorInput from '../AssertionFormSelectorInput';

export default {
  title: 'Assertion Form Selector Input',
  component: AssertionFormSelectorInput,
  argTypes: {onChange: {action: 'onChange'}},
} as ComponentMeta<typeof AssertionFormSelectorInput>;

const Template: ComponentStory<typeof AssertionFormSelectorInput> = args => <AssertionFormSelectorInput {...args} />;

const span = SpanMock.model();

export const Default = Template.bind({});
Default.args = {
  attributeList: span.attributeList,
  value: [],
};
