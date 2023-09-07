# Creating Tests

![Main Screen](../img/main-screen-0.11.png)

Click the **Create** button and select **Create New Test** in the drop down:

![Create a Test Button](../img/create-button-0.11.png)

The "Create New Test" dialog appears:

![Create a Test](../img/create-test-0.13.png)

The option to choose the kind of trigger to initiate the trace is presented:

- HTTP Request - Create a basic HTTP request.
- GRPC Request - Test and debug your GRPC request.
- cURL Command - Define your HTTP test via a cURL command.
- Postman Collection - Define your HTTP request via a Postman collection.
- TraceID - Define your test via a TraceID.
- Kafka - Test consumers with Kafka messages

Choose the trigger and click **Next**:

![Choose Trigger](../img/choose-trigger-0.13.png)

In this example, HTTP Request has been chosen.

![Choose Example](../img/choose-example-0.11.png)

Input the **Name** of the test and the **Description** or select one of the example provided in the drop down:

![Choose Example Pokemon](../img/choose-example-pokemon-0.11.png)

The **Pokemon - Import** example has been chosen. Then click **Next**.

![Choose Example Pokemon](../img/choose-example-pokemon-import-0.11.png)

Add any additional information and click **Create & Run**:

![Create Test](../img/provide-addl-information-0.11.png)

The test will start:

![Awaiting Trace](../img/awaiting-trace-0.11.png)

When the test is finished, you will get the following results:

![Finished Trace](../img/finished-trace-0.11.png)

Please visit the [Test Results](test-results.md) document for an explanation of viewing the results of a test.
