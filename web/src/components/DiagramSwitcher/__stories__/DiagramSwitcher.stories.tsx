import {ComponentStory, ComponentMeta} from '@storybook/react';

import DiagramSwitcher from '../DiagramSwitcher';

export default {
  title: 'Diagram Switcher',
  component: DiagramSwitcher,
  argTypes: {onClick: {action: 'clicked'}, onDelete: {action: 'deleted'}, onRunTest: {action: 'runTest'}},
} as ComponentMeta<typeof DiagramSwitcher>;

const Template: ComponentStory<typeof DiagramSwitcher> = args => <DiagramSwitcher {...args} />;

export const Default = Template.bind({});
Default.args = {};
