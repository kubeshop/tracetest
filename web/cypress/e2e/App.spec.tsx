import {mount} from '@cypress/react';
import App from '../../src/App';

it('renders learn react link', () => {
  mount(<App />);
  cy.get('h3.ant-typography').contains(/all tests/i);
});
