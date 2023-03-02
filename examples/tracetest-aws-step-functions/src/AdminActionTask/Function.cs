using System;
using Amazon.Lambda.Core;
using Amazon.XRay.Recorder.Handlers.AwsSdk;
using Plagiarism;
using PlagiarismRepository;

// Assembly attribute to enable the Lambda function's JSON input to be converted into a .NET class.
[assembly: LambdaSerializer(typeof(Amazon.Lambda.Serialization.SystemTextJson.DefaultLambdaJsonSerializer))]
namespace AdminActionTask
{
  public class Function
  {

    private readonly IIncidentRepository _incidentRepository;

    public Function()
    {
      var tableName = Environment.GetEnvironmentVariable("TABLE_NAME");
      _incidentRepository = new IncidentRepository(tableName);
    }

    public Function(IIncidentRepository incidentRepository)
    {
      _incidentRepository = incidentRepository;
      AWSSDKHandler.RegisterXRayForAllServices();
    }

    /// <summary>
    /// A simple function that takes a string and does a ToUpper
    /// </summary>
    /// <param name="incident"></param>
    /// <param name="context"></param>
    /// <returns></returns>
    public void FunctionHandler(Incident incident, ILambdaContext context)
    {
      incident.AdminActionRequired = true;
      incident.IncidentResolved = false;
      incident.ResolutionDate = DateTime.Now;

      _incidentRepository.SaveIncident(incident);
    }
  }
}
