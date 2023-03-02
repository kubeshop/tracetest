using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Amazon;
using Amazon.DynamoDBv2;
using Amazon.DynamoDBv2.Model;
using Amazon.Lambda.APIGatewayEvents;
using Amazon.Lambda.TestUtilities;
using Amazon.StepFunctions;
using Newtonsoft.Json;
using Xunit;
using Xunit.Abstractions;

namespace SubmitExamTask.Tests
{
    public class FunctionTests : IDisposable
    {
        private readonly ITestOutputHelper _testOutputHelper;
        private readonly IAmazonDynamoDB _dynamoDbClient;
        private const string TablePrefix = "IncidentsTestTable-";

        private const string Token =
            "AAAAKgAAAAIAAAAAAAAAAbdvA5UnsPbXk2HGkayUMygJK8eFJq3pnwBV/xTTDwiIbXvk246zL6Y1+UxXRWzPnbLD0mex2AEUEwMfjxjOj0lW0g+6AwFv6gA0MW/gU2SAdkHZl7tQQ1o3uBL2eOlSSYakcvPvF35BJdXCFkhhKaoqB8CzpnzkJPr7KVSXumjMouy/C4KwJJMqcVpeIW2Xhjyxq6FFT8+GRfNspJUaGE3aId15q/dK94xRTPG/Gidez7iuINk6Y7JpbA4/sj3T2hpUuDKyi4CcCkI8A4z93Hn2Tw2OMqWwhmserDGNfI3UgW3Um6pHRYNvL1prARZ9DkGHHftGaaXXBU8IO1mxYij4TciyP2Cky4b/Dk6ImioM0s+xdIeFOfMprMg73KG5WPK0XAWF+coMC7zBKJTtHZmudk9wKzTPdiSEZrwmPgeD3hVeWTQXwi7GF9hVbpS8wz/QrtI78HGPcbUdMi0Y79YihuGDo6iN4booO/5Tek3prcfDKhU3JtqqqVFRp9ugqQlOxhnkGmKaajp5miRFDcgrghxvP8Fp4D1DDY+/5vUxHFS+tOqvrp24YpSfO51xQxp7GWeg0k9qSnSWntOKdJRjmE7gyvIhKC9XMnlLktJEeBpCQa/B3pqzIr31sPB9ooDTS7m97REIl6Gf0VOtOx4=";
        private string _tableName;


        public FunctionTests(ITestOutputHelper testOutputHelper)
        {
            _testOutputHelper = testOutputHelper;
            _dynamoDbClient = new AmazonDynamoDBClient(RegionEndpoint.APSoutheast2);

            SetupTableAsync().Wait();
            
            
        }


        [Fact]
        public void SubmitExamHandlerReturns200ForValidRequest()
        {
            var request = new APIGatewayProxyRequest();
            var body = new Dictionary<string, string>
            {
                {"IncidentId", "d557cc5a-4bcc-43ec-914b-4498896aedc4"},
                {"ExamId", "b439492b-db9a-4ec7-b5c5-3da29d37d874"},
                {"Score", "99"},
                {"TaskToken", Token}
            };
            request.Body = JsonConvert.SerializeObject(body);

            var context = new TestLambdaContext();

            var expectedResponse = new APIGatewayProxyResponse
            {
                StatusCode = 200,
                Headers = new Dictionary<string, string> {{"Content-Type", "application/json"}}
            };

            var function = new Function(_dynamoDbClient, new AmazonStepFunctionsClient(RegionEndpoint.APSoutheast2),
                _tableName);
            var response = function.FunctionHandler(request, context);

            _testOutputHelper.WriteLine("Lambda Response: \n" + response.StatusCode);
            _testOutputHelper.WriteLine("Expected Response: \n" + expectedResponse.StatusCode);

            Assert.Equal(expectedResponse.Body, response.Body);
            Assert.Equal(expectedResponse.Headers, response.Headers);
            Assert.Equal(expectedResponse.StatusCode, response.StatusCode);
        }

        [Fact]
        public void SubmitExamHandlerReturns400ForInvalidExamIdRequest()
        {
            var request = new APIGatewayProxyRequest();

            var body = new Dictionary<string, string>
            {
                {"ExamId", "dc-not-a-guid"},
                {"IncidentId", "6b44fd97-1af3-42f6-9a0b-0138fffa8cf4"},
                {"Score", "65"},
                {"TaskToken", Token}
            };
            request.Body = JsonConvert.SerializeObject(body);

            var context = new TestLambdaContext();

            var expectedResponse = new APIGatewayProxyResponse
            {
                StatusCode = 400,
                Headers = new Dictionary<string, string> {{"Content-Type", "application/json"}}
            };

            var function = new Function(_dynamoDbClient, new AmazonStepFunctionsClient(RegionEndpoint.APSoutheast2),
                _tableName);
            var response = function.FunctionHandler(request, context);

            _testOutputHelper.WriteLine("Lambda Response: \n" + response.StatusCode);
            _testOutputHelper.WriteLine("Expected Response: \n" + expectedResponse.StatusCode);

            Assert.Equal(expectedResponse.Headers, response.Headers);
            Assert.Equal(expectedResponse.StatusCode, response.StatusCode);
        }

        [Fact]
        public void SubmitExamHandlerReturns400ForInvalidIncidentIdRequest()
        {
            var request = new APIGatewayProxyRequest();

            var body = new Dictionary<string, string>
            {
                {"ExamId", "dc149d4b-ce6d-435a-b922-9da90f7c3eed"},
                {"IncidentId", "not-a-guid-0138fffa8cf4"},
                {"Score", "65"},
                {"TaskToken", Token}
            };
            request.Body = JsonConvert.SerializeObject(body);

            var context = new TestLambdaContext();

            var expectedResponse = new APIGatewayProxyResponse
            {
                StatusCode = 400,
                Headers = new Dictionary<string, string> {{"Content-Type", "application/json"}}
            };

            var function = new Function(_dynamoDbClient, new AmazonStepFunctionsClient(RegionEndpoint.APSoutheast2),
                _tableName);
            var response = function.FunctionHandler(request, context);

            _testOutputHelper.WriteLine("Lambda Response: \n" + response.StatusCode);
            _testOutputHelper.WriteLine("Expected Response: \n" + expectedResponse.StatusCode);

            Assert.Equal(expectedResponse.Headers, response.Headers);
            Assert.Equal(expectedResponse.StatusCode, response.StatusCode);
        }

        /// <summary>
        /// Helper function to create a testing table
        /// </summary>
        /// <returns></returns>
        private async Task SetupTableAsync()
        {
            var listTablesResponse = await _dynamoDbClient.ListTablesAsync(new ListTablesRequest());
            var existingTestTable =
                listTablesResponse.TableNames.FindAll(s => s.StartsWith(TablePrefix)).FirstOrDefault();
            if (existingTestTable == null)
            {
                _tableName = TablePrefix + DateTime.Now.Ticks;

                CreateTableRequest request = new CreateTableRequest
                {
                    TableName = _tableName,
                    ProvisionedThroughput = new ProvisionedThroughput
                    {
                        ReadCapacityUnits = 2,
                        WriteCapacityUnits = 2
                    },
                    KeySchema = new List<KeySchemaElement>
                    {
                        new KeySchemaElement
                        {
                            AttributeName = "IncidentId"
                        }
                    },
                    AttributeDefinitions = new List<AttributeDefinition>
                    {
                        new AttributeDefinition
                        {
                            AttributeName = "IncidentId",
                            AttributeType = ScalarAttributeType.S
                        }
                    }
                };

                await _dynamoDbClient.CreateTableAsync(request);

                var describeRequest = new DescribeTableRequest {TableName = _tableName};
                DescribeTableResponse response;

                do
                {
                    Thread.Sleep(1000);
                    response = await _dynamoDbClient.DescribeTableAsync(describeRequest);
                } while (response.Table.TableStatus != TableStatus.ACTIVE);
            }
            else
            {
                Console.WriteLine($"Using existing test table {existingTestTable}");
                _tableName = existingTestTable;
            }
        }

        public void Dispose()
        {
            _dynamoDbClient?.Dispose();
        }
    }
}