using System;
using System.Collections.Generic;
using System.Net;
using Amazon.DynamoDBv2;
using Amazon.Lambda.APIGatewayEvents;
using Amazon.Lambda.Core;
using Amazon.StepFunctions;
using Amazon.StepFunctions.Model;
using Amazon.XRay.Recorder.Handlers.AwsSdk;
using Newtonsoft.Json;
using PlagiarismRepository;

// Assembly attribute to enable the Lambda function's JSON input to be converted into a .NET class.
[assembly: LambdaSerializer(typeof(Amazon.Lambda.Serialization.SystemTextJson.DefaultLambdaJsonSerializer))]

namespace SubmitExamTask
{
    public class Function
    {
        private readonly IIncidentRepository _incidentRepository;
        private readonly AmazonStepFunctionsClient _amazonStepFunctionsClient;

        /// <summary>
        /// Default constructor
        /// </summary>
        public Function()
        {
            AWSSDKHandler.RegisterXRayForAllServices();
            _incidentRepository = new IncidentRepository(Environment.GetEnvironmentVariable("TABLE_NAME"));
            _amazonStepFunctionsClient = new AmazonStepFunctionsClient();
        }

        /// <summary>
        /// Constructor used for testing purposes
        /// </summary>
        /// <param name="ddbClient"></param>
        /// <param name="stepFunctions"></param>
        /// <param name="tablename"></param>
        public Function(IAmazonDynamoDB ddbClient, IAmazonStepFunctions stepFunctions, string tablename)
        {
            AWSSDKHandler.RegisterXRayForAllServices();
            _incidentRepository = new IncidentRepository(ddbClient, tablename);
            _amazonStepFunctionsClient = (AmazonStepFunctionsClient) stepFunctions;
        }

        /// <summary>
        /// A simple function that takes a string and does a ToUpper
        /// </summary>
        /// <param name="request">Instance of APIGatewayProxyRequest</param>
        /// <param name="context">AWS Lambda Context</param>
        /// <returns>Instance of APIGatewayProxyResponse</returns>
        public APIGatewayProxyResponse FunctionHandler(APIGatewayProxyRequest request, ILambdaContext context)
        {
            var body = JsonConvert.DeserializeObject<Dictionary<string, string>>(request?.Body);

            var isIncidentId = Guid.TryParse(body["IncidentId"], out var incidentId);
            var isExamId = Guid.TryParse(body["ExamId"], out var examId);    
            var isScore = Int32.TryParse(body["Score"], out var score);
            var token = body["TaskToken"];

            if (!isIncidentId || !isExamId | !isScore | !(token.Length >= 1 & token.Length <= 1024))
            {
                return new APIGatewayProxyResponse
                {
                    StatusCode = (int) HttpStatusCode.BadRequest,
                    Headers = new Dictionary<string, string> {{"Content-Type", "application/json"}}
                };
            }

            Console.WriteLine($"IncidentId: {incidentId}");
            Console.WriteLine($"ExamId: {examId}");
            Console.WriteLine($"Score: {score}");
            Console.WriteLine($"Token: {token}");

            var incident = _incidentRepository.GetIncidentById(incidentId);
            var exam = incident.Exams.Find(e => e.ExamId == examId);
            exam.Score = score;

            _incidentRepository.SaveIncident(incident);

            Console.WriteLine(JsonConvert.SerializeObject(incident));

            var sendTaskSuccessRequest = new SendTaskSuccessRequest
            {
                TaskToken = token,
                Output = JsonConvert.SerializeObject(incident)
            };
            
            try
            {
                _amazonStepFunctionsClient.SendTaskSuccessAsync(sendTaskSuccessRequest).Wait();
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
                throw;
            }

            return new APIGatewayProxyResponse
            {
                StatusCode = (int) HttpStatusCode.OK,
                Headers = new Dictionary<string, string> {
                    {"Content-Type", "application/json"}, 
                    {"Access-Control-Allow-Origin", "*"},
                    {"Access-Control-Allow-Headers", "Content-Type"},
                    {"Access-Control-Allow-Methods", "OPTIONS,POST"}
                }
            };
        }
    }
}