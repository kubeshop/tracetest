const Tracetest = require('@tracetest/client').default;

const { TRACETEST_API_TOKEN = '' } = process.env;

async function runTracebasedTest(testDefinition, traceID) {
  const tracetestClient = await Tracetest({ apiToken: TRACETEST_API_TOKEN });

  const test = await tracetestClient.newTest(testDefinition);
  await tracetestClient.runTest(test, { variables: [ { key: 'TRACE_ID', value: traceID }] });
  console.log(await tracetestClient.getSummary());
}

module.exports = { runTracebasedTest };
