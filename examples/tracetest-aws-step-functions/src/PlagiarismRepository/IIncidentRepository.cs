using System;
using Plagiarism;

namespace PlagiarismRepository
{
    public interface IIncidentRepository
    {
        Incident SaveIncident(Incident incident);
        Incident GetIncidentById(Guid incidentId);
    }
}