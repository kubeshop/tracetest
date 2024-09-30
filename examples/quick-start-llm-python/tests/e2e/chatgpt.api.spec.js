// @ts-check
const { test, expect } = require('@playwright/test');

const chatgptTraceBasedTest = require('./definitions/chatgpt');

const { runTracebasedTest } = require('./tracetest');

test('generated summarized test for OpenAI', async ({ request }) => {
  const result = await request.post(`http://localhost:8800/summarizeText`, {
    data: {
      provider: "OpenAI (ChatGPT)",
      text: "Born in London, Turing was raised in southern England. He graduated from King's College, Cambridge, and in 1938, earned a doctorate degree from Princeton University. During World War II, Turing worked for the Government Code and Cypher School at Bletchley Park, Britain's codebreaking centre that produced Ultra intelligence. He led Hut 8, the section responsible for German naval cryptanalysis. Turing devised techniques for speeding the breaking of German ciphers, including improvements to the pre-war Polish bomba method, an electromechanical machine that could find settings for the Enigma machine. He played a crucial role in cracking intercepted messages that enabled the Allies to defeat the Axis powers in many crucial engagements, including the Battle of the Atlantic.\n\nAfter the war, Turing worked at the National Physical Laboratory, where he designed the Automatic Computing Engine, one of the first designs for a stored-program computer. In 1948, Turing joined Max Newman's Computing Machine Laboratory at the Victoria University of Manchester, where he helped develop the Manchester computers[12] and became interested in mathematical biology. Turing wrote on the chemical basis of morphogenesis and predicted oscillating chemical reactions such as the Belousov–Zhabotinsky reaction, first observed in the 1960s. Despite these accomplishments, he was never fully recognised during his lifetime because much of his work was covered by the Official Secrets Act."
    }
  });

  const jsonResult = await result.json();
  expect(jsonResult).not.toBe(null);
  expect(jsonResult.summary).not.toBe(null);

  const traceID = jsonResult.trace_id;
  expect(traceID).not.toBe(null);

  // run trace-based test
  await runTracebasedTest(chatgptTraceBasedTest, traceID);
});
