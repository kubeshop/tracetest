using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Amazon;
using Amazon.DynamoDBv2;
using Amazon.DynamoDBv2.Model;
using Xunit;
using Amazon.Lambda.TestUtilities;
using NSubstitute;
using Plagiarism;
using PlagiarismRepository;

namespace ResolveIncidentTask.Tests
{
    public class FunctionTests
    {
        private readonly IAmazonDynamoDB _dynamoDbClient;
        private const string TablePrefix = "IncidentsTestTable-";
        private string _tableName;
        
        public FunctionTests()
        {
            _dynamoDbClient = new AmazonDynamoDBClient(RegionEndpoint.APSoutheast2);
        }

        [Fact]
        public void ResolveIncidentFunctionTest()
        {
            var incidentRepository
                = Substitute.For<IIncidentRepository>();


            var function = new Function(_dynamoDbClient, _tableName);
            var context = new TestLambdaContext();

            var state = new Incident
            {
                IncidentId = Guid.NewGuid(),
                StudentId = "123",
                IncidentDate = new DateTime(2018, 02, 03),
                Exams = new List<Exam>()
                {
                    new Exam(Guid.NewGuid(), new DateTime(2018, 02, 10), 10),
                    new Exam(Guid.NewGuid(), new DateTime(2018, 02, 17), 65)
                },
                ResolutionDate = null
            };


            incidentRepository.SaveIncident(state);
            incidentRepository.ReceivedCalls();

            function.FunctionHandler(state, context);
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
    }
}