import {ComponentStory, ComponentMeta} from '@storybook/react';
import {Plugins} from 'constants/Plugins.constants';

import PluginCard from '../PluginCard';

export default {
  title: 'Plugin Card',
  component: PluginCard,
  argTypes: {onSelect: {action: 'onSelect'}},
} as ComponentMeta<typeof PluginCard>;

const Template: ComponentStory<typeof PluginCard> = args => <PluginCard {...args} />;

export const Default = Template.bind({});
Default.args = {
  plugin: Plugins.REST,
};
