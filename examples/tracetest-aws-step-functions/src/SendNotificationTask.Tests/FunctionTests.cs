using System;
using System.Collections.Generic;
using Amazon.Lambda.TestUtilities;
using Plagiarism;
using Xunit;

namespace SendNotificationTask.Tests
{
    public class FunctionTests
    {
        private readonly TestLambdaContext _context;
        private readonly IncidentWrapper _incidentIn;

        public FunctionTests()
        {
            _context = new TestLambdaContext();

            _incidentIn = new IncidentWrapper()
            {
                Input = new Incident
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
                },
                TaskToken = "TASKTOKEN"
            };
        }


        [Fact]
        public void NotificationSentSouldBeFalseIfSnsPublishSucceeds()
        {
            var function = new Function();
            function.FunctionHandler(_incidentIn, _context);

            //Assert.True(response.Exams[0].NotificationSent == true);
        }

        [Fact]
        public void NotificationSentSouldBeFalseIfSnsPublishFails()
        {

            var function = new Function();

            function.FunctionHandler(_incidentIn, _context);

             // Assert.True(response.Exams[0].NotificationSent == false);
        }
    }
}