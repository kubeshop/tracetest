import {ComponentStory, ComponentMeta} from '@storybook/react';

import FileViewerModal from '../FileViewerModal';

export default {
  title: 'File Viewer Modal',
  component: FileViewerModal,
  argTypes: {onClose: {action: 'onClose'}},
} as ComponentMeta<typeof FileViewerModal>;

const Template: ComponentStory<typeof FileViewerModal> = args => <FileViewerModal {...args} />;

export const TestDefinition = Template.bind({});
TestDefinition.args = {
  data: `name: POST import pokemon
  description: Import a pokemon using its ID
  trigger:
    type: http
    http_request:
      url: http://demo-pokemon-api.demo/pokemon/import
      method: POST
      headers:
      - key: Content-Type
        value: application/json
      body:
        type: raw
        raw: '{ "id": 52 }'
  testDefinition:
  - selector: span[name = "POST /pokemon/import"]
    assertions:
    - tracetest.span.duration <= 500
    - http.status_code = 200
  - selector: span[name = "send message to queue"]
    assertions:
    - messaging.message.payload contains 52
  - selector: span[name = "consume message from queue"]:last
    assertions:
    - messaging.message.payload contains 52
  - selector: span[name = "consume message from queue"]:last span[name = "import pokemon
      from pokeapi"]
    assertions:
    - http.status_code = 200
  - selector: span[name = "consume message from queue"]:last span[name = "save pokemon
      on database"]
    assertions:
    - db.repository.operation = "create"
    - tracetest.span.duration <= 500
`,
  isOpen: true,
  title: 'Test Definition',
  language: 'yaml',
  subtitle: 'Preview your YAML file',
  fileName: 'definition.yaml',
};

export const JUnit = Template.bind({});
JUnit.args = {
  data: `<testsuites name="POST import pokemon" tests="6" failures="3" errors="0" skipped="0" time="8">
	<testsuite name="span[name = &#34;POST /pokemon/import&#34;]" tests="2" failures="0" errors="0" skipped="0">
		<testcase name="&#34;tracetest.span.duration&#34; &lt;= &#34;500&#34;"></testcase>
		<testcase name="&#34;http.status_code&#34; = &#34;200&#34;"></testcase>
	</testsuite>
	<testsuite name="span[name = &#34;send message to queue&#34;]" tests="1" failures="1" errors="0" skipped="0">
		<testcase name="&#34;messaging.message.payload&#34; contains &#34;52&#34;">
			<failure type="messaging.message.payload" message=""></failure>
		</testcase>
	</testsuite>
	<testsuite name="span[name = &#34;consume message from queue&#34;]:last" tests="1" failures="1" errors="0" skipped="0">
		<testcase name="&#34;messaging.message.payload&#34; contains &#34;52&#34;">
			<failure type="messaging.message.payload" message=""></failure>
		</testcase>
	</testsuite>
	<testsuite name="span[name = &#34;consume message from queue&#34;]:last span[name = &#34;import pokemon from pokeapi&#34;]" tests="0" failures="0" errors="0" skipped="0"></testsuite>
	<testsuite name="span[name = &#34;consume message from queue&#34;]:last span[name = &#34;save pokemon on database&#34;]" tests="2" failures="1" errors="0" skipped="0">
		<testcase name="&#34;db.repository.operation&#34; = &#34;create&#34;">
			<failure type="db.repository.operation" message=""></failure>
		</testcase>
		<testcase name="&#34;tracetest.span.duration&#34; &lt;= &#34;500&#34;"></testcase>
	</testsuite>
</testsuites>
`,
  isOpen: true,
  title: 'JUnit Results',
  language: 'xml',
  subtitle: 'Preview your JUnit results',
  fileName: 'junit.xml',
};
