using System;
using Amazon.DynamoDBv2;
using Amazon.Lambda.Core;
using Amazon.XRay.Recorder.Handlers.AwsSdk;
using Plagiarism;
using PlagiarismRepository;

// Assembly attribute to enable the Lambda function's JSON input to be converted into a .NET class.
[assembly: LambdaSerializer(typeof(Amazon.Lambda.Serialization.SystemTextJson.DefaultLambdaJsonSerializer))]
namespace ResolveIncidentTask
{
    public class Function
    {
        private readonly IIncidentRepository _incidentRepository;

        public Function()
        {
            AWSSDKHandler.RegisterXRayForAllServices();
            _incidentRepository = new IncidentRepository(Environment.GetEnvironmentVariable("TABLE_NAME"));
        }

        /// <summary>
        /// Constructor used for testing purposes
        /// </summary>
        /// <param name="ddbClient">Instance of DynamoDB client</param>
        /// <param name="tablename">DynamoDB table name</param>
        public Function(IAmazonDynamoDB ddbClient, string tablename)
        {
            AWSSDKHandler.RegisterXRayForAllServices();
            _incidentRepository = new IncidentRepository(ddbClient, tablename);
        }

        /// <summary>
        /// Function to resolve the incident and cpmplete the workflow.
        /// All state data is persisted.
        /// </summary>
        /// <param name="incident"></param>
        /// <param name="context"></param>
        /// <returns></returns>
        public void FunctionHandler(Incident incident, ILambdaContext context)
        {
            incident.AdminActionRequired = false;
            incident.IncidentResolved = true;
            incident.ResolutionDate = DateTime.Now;

            _incidentRepository.SaveIncident(incident);
        }
    }


    
}