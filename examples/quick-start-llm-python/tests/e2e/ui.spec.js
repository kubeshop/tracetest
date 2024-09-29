// @ts-check
const { test, expect } = require('@playwright/test');

const geminiTraceBasedTest = require('./definitions/gemini');
const chatgptTraceBasedTest = require('./definitions/chatgpt');

const { runTracebasedTest } = require('./tracetest');

const timeToWait = 10_000; // 10 seconds

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

test('generated summarized test for Gemini', async ({ page }) => {
  // Go to Streamlit app
  await page.goto('http://localhost:8501/');

  // Select Google (Gemini) model
  await page.getByTestId('stSelectbox').locator('div').filter({ hasText: 'Google (Gemini)' }).nth(2).click();

  // Click on add example text
  await page.getByRole('button', { name: 'Add example text' }).click();

  // Click on button to call summarization rule
  await page.getByRole('button', { name: 'Summarize' }).click();

  // Wait for time
  await sleep(timeToWait);

  // Capture TraceID
  const traceIDLabel = await page.getByRole('link', { name: 'Trace ID' });
  expect(traceIDLabel).toHaveText('Trace ID');

  console.log(traceIDLabel);

  // const traceID = (traceIDLabel || '').replace('Trace ID:', '').trim();

  // // run trace-based test
  // await runTracebasedTest(geminiTraceBasedTest, traceID);
});

// test('generated summarized test for OpenAPI', async ({ page }) => {
//   // Go to Streamlit app
//   await page.goto('http://localhost:8501/');

//   // Select OpenAI (ChatGPT) model
//   await page.getByTestId('stSelectbox').locator('div').filter({ hasText: 'OpenAI (ChatGPT)' }).nth(2).click();

//   // Click on add example text
//   await page.getByRole('button', { name: 'Add example text' }).click();

//   // Click on button to call summarization rule
//   await page.getByRole('button', { name: 'Summarize' }).click();

//   // Wait for time
//   await sleep(timeToWait);

//   // Capture TraceID
//   const traceIDElement = await page.getByText('Trace ID:');
//   expect(traceIDElement).toHaveText('Trace ID:');

//   const traceIDLabel = await page.getByText('Trace ID:').innerText();
//   expect(traceIDLabel).not.toBeNull();

//   const traceID = (traceIDLabel || '').replace('Trace ID:', '').trim();

//   // run trace-based test
//   await runTracebasedTest(chatgptTraceBasedTest, traceID);
// });
