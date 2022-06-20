import {renderHook} from '@testing-library/react-hooks';
import {Button, Form} from 'antd';
import {render} from 'test-utils';
import {clickButton, selectInput, typeInputValue} from 'utils/Tests';
import CreateTestForm, {ICreateTestValues} from '../CreateTestForm';

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

  it('should call onSubmit', async () => {
    const onSubmit = jest.fn();
    const {
      result: {
        current: [form],
      },
    } = renderHook(() => Form.useForm<ICreateTestValues>());
    render(
      <>
        <CreateTestForm
          form={form}
          onSelectDemo={jest.fn()}
          onSubmit={onSubmit}
          onValidation={jest.fn()}
          selectedDemo=""
        />
        <Button data-cy="create-test-submit" onClick={() => form.submit()}>
          Create
        </Button>
      </>
    );
    await typeInputValue('name', 'My cool test', 'input');
    await selectInput('method-select', 'method-select-option-GET');
    await typeInputValue('url', 'https://google.com');
    await typeInputValue('body', '{id: 3}', 'textarea');
    await selectInput('auth-type-select', 'auth-type-select-option-apiKey');
    await typeInputValue('url', 'https://google.com');
    await typeInputValue('apiKey-key', 'key', 'input');
    await typeInputValue('apiKey-value', 'value', 'input');
    await selectInput('auth-apiKey-select', 'auth-apiKey-select-option-header');
    await clickButton('create-test-submit');
    expect(onSubmit).toBeCalled();
  });
});
