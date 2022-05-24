import {render} from 'test-utils';

import CreateTestModal from '../CreateTestModal';

describe('CreateTestModal', () => {
  it('should render', () => {
    const {getAllByText, getByPlaceholderText, getByText} = render(<CreateTestModal onClose={jest.fn()} visible />);

    expect(getByPlaceholderText('Enter request url')).toBeInTheDocument();
    expect(getByPlaceholderText('Enter test name')).toBeInTheDocument();
    expect(getByPlaceholderText('Enter request body text')).toBeInTheDocument();
    expect(getByText(/cancel/i)).toBeInTheDocument();
    expect(getAllByText(/create/i).length).toBe(2);
  });
});
