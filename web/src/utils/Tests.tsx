import {act, fireEvent, screen} from '@testing-library/react';

export async function typeInputValue(key: string, value: string, input?: 'input' | 'textarea') {
  const byTestId = screen.getByTestId(key);
  const querySelector = input ? byTestId.querySelector(input) : byTestId;
  await act(async () => {
    fireEvent.change(querySelector as Element, {target: {value}});
  });
  expect(querySelector).toHaveValue(value);
}

export async function selectInput(selectorTestId: string, optionTestId: string) {
  const byTestId = screen.getByTestId(selectorTestId).querySelector('input');
  await act(async () => {
    fireEvent.click(byTestId as Element);
    fireEvent.mouseDown(byTestId as Element);
  });
  const selection = screen.getByTestId(optionTestId);
  await act(async () => {
    fireEvent.click(selection as Element);
  });
}

export async function clickButton(testId: string) {
  const byTestId = screen.getByTestId(testId);
  await act(async () => {
    fireEvent.click(byTestId as Element);
  });
}
