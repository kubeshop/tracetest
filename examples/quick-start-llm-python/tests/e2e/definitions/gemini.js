const definition = {
  "type": "Test",
  "spec": {
    "id": "VS0U-HgHg",
    "name": "Get Gemini trace",
    "trigger": {
      "type": "traceid",
      "traceid": {
        "id": "${var:TRACE_ID}"
      }
    },
    "specs": [
      {
        "selector": "span[tracetest.span.type=\"general\" name=\"MapReduceDocumentsChain.workflow\"]",
        "name": "It triggered a Summarization workflow",
        "assertions": [
          "attr:traceloop.workflow.name   =   \"MapReduceDocumentsChain\""
        ]
      },
      {
        "selector": "span[tracetest.span.type=\"general\" name=\"ChatGoogleGenerativeAI.chat\"]",
        "name": "It called Gemini API at least once",
        "assertions": [
          "attr:tracetest.selected_spans.count   >=   1"
        ]
      }
    ],
    "pollingProfile": "predefined-default"
  }
};

module.exports = definition;
