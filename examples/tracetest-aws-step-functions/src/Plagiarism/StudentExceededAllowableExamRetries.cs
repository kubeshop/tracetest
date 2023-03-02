using System;

namespace Plagiarism
{
    /// <summary>
    /// Custom Exception for students that have exceeded the allowable number of exam retries.
    /// </summary>
    /// <seealso cref="T:System.Exception" />
    public class StudentExceededAllowableExamRetries : Exception
    {
        public StudentExceededAllowableExamRetries()
        {
        }

        public StudentExceededAllowableExamRetries(string message) : base((string) message)
        {
        }
    }
}
