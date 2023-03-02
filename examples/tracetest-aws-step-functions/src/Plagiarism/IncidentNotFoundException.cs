using System;

namespace Plagiarism
{
    public class IncidentNotFoundException : Exception
    {
        public IncidentNotFoundException()
        { 
        }
        public IncidentNotFoundException(string message) : base((string) message)
        {
        }
    }
}