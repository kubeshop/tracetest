const definition = {
  "type": "Test",
  "spec": {
    "id": "B9opfNRNR",
    "name": "Get GPT4 trace",
    "trigger": {
      "type": "traceid",
      "traceid": {
        "id": "${var:TRACE_ID}"
      }
    },
    "specs": [
      {
        "selector": "span[tracetest.span.type=\"general\" name=\"ChatPromptTemplate.workflow\"]",
        "name": "It performed a Chat workflow",
        "assertions": [
          "attr:tracetest.span.name = \"ChatPromptTemplate.workflow\""
        ]
      },
      {
        "selector": "span[tracetest.span.type=\"general\" name=\"openai.chat\"]",
        "name": "It called OpenAI API",
        "assertions": [
          "attr:name = \"openai.chat\""
        ]
      }
    ],
    "pollingProfile": "predefined-default"
  }
};

module.exports = definition;
