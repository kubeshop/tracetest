import faker from '@faker-js/faker';
import {ComponentStory, ComponentMeta} from '@storybook/react';

import AttributeRow from '../AttributeRow';

export default {
  title: 'Attribute Row',
  component: AttributeRow,
  argTypes: {onCreateAssertion: {action: 'clicked'}},
} as ComponentMeta<typeof AttributeRow>;

const Template: ComponentStory<typeof AttributeRow> = args => (
  <div style={{width: 670}}>
    <AttributeRow {...args} />
  </div>
);

export const Default = Template.bind({});
Default.args = {
  attribute: {
    key: faker.random.word(),
    value:
      '{"totalCount":327,"items":[{"id":1,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png","isFeatured":false,"type":"electric","name":"pikachu"},{"id":2,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png","isFeatured":false,"type":"electric","name":"pikachu"},{"id":3,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png","isFeatured":false,"type":"electric","name":"pikachu"},{"id":4,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png","isFeatured":false,"type":"electric","name":"pikachu"},{"id":5,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png","isFeatured":false,"type":"electric","name":"pikachu"},{"id":6,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":7,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":8,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":9,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":10,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":11,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":12,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":13,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":14,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":15,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":16,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":17,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":18,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":19,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"},{"id":20,"imageUrl":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png","isFeatured":false,"type":"normal","name":"meowth"}]}',
  },
};
