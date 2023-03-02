using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Amazon;
using Amazon.DynamoDBv2;
using Amazon.DynamoDBv2.Model;
using Plagiarism;
using Xunit;

namespace PlagiarismRepository.Tests
{
    public class RepositoryTests
    {
        private const string TablePrefix = "IncidentsTestTable-";
        private string _tableName;
        private readonly IAmazonDynamoDB _dynamoDbClient;
        private IIncidentRepository _incidentRepository;


        public RepositoryTests()
        {
            _dynamoDbClient = new AmazonDynamoDBClient(RegionEndpoint.APSoutheast2);
            SetupTableAsync().Wait();
        }

        [Fact]
        public void SaveIncidentAsync()
        {
            _incidentRepository = new IncidentRepository(_dynamoDbClient, _tableName);

            var state = new Incident
            {
                IncidentId = Guid.NewGuid(),
                StudentId = "123",
                IncidentDate = new DateTime(2018, 02, 03),
                Exams = new List<Exam>()
                {
                    new Exam(Guid.NewGuid(), new DateTime(2018, 02, 17), 0),
                    new Exam(Guid.NewGuid(), new DateTime(2018, 02, 10), 65)
                },
                ResolutionDate = null
            };

            var incident =  _incidentRepository.SaveIncident(state);

            Assert.NotNull(incident);
        }

        [Fact]
        public void UpdateIncidentAsync()
        {
            _incidentRepository = new IncidentRepository(_dynamoDbClient, _tableName);

            var state = new Incident
            {
                IncidentId = Guid.NewGuid(),
                StudentId = "123",
                IncidentDate = new DateTime(2018, 02, 03),
            };

            var incident =  _incidentRepository.SaveIncident(state);

            incident.Exams = new List<Exam>
            {
                new Exam(Guid.NewGuid(), new DateTime(2018, 02, 17), 0),
                new Exam(Guid.NewGuid(), new DateTime(2018, 02, 10), 65),
                new Exam(Guid.NewGuid(), new DateTime(2018, 02, 17), 99)
            };

            var updatedIncident = _incidentRepository.SaveIncident(state);

            Assert.NotNull(incident);
            Assert.True(updatedIncident.Exams.Count == 3, "Should be three");
        }


        [Fact]
        public void FindIncidentAsync()
        {
            _incidentRepository = new IncidentRepository(_dynamoDbClient, _tableName);

            var state = new Incident
            {
                IncidentId = Guid.NewGuid(),
                StudentId = "123",
                IncidentDate = new DateTime(2018, 02, 03),
                Exams = new List<Exam>()
                {
                    new Exam(Guid.NewGuid(), new DateTime(2018, 02, 17), 0),
                    new Exam(Guid.NewGuid(), new DateTime(2018, 02, 10), 65)
                },
                ResolutionDate = null
            };

            var incident =  _incidentRepository.SaveIncident(state);

            var newIncident =  _incidentRepository.GetIncidentById(incident.IncidentId);

            Assert.NotNull(newIncident);
            Assert.True(newIncident.IncidentId == incident.IncidentId, "Should be teh same incident");
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