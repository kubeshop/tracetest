using System;
using System.Collections.Generic;

namespace Plagiarism
{
    public class IncidentWrapper
    {
        public Incident Input { get; set; }
        public string TaskToken { get; set; }
    }

    public class Incident
    {
        public string StudentId { get; set; }
        public Guid IncidentId { get; set; }
        public DateTime IncidentDate { get; set; }
        public List<Exam> Exams { get; set; }
        public DateTime? ResolutionDate { get; set; }
        public bool IncidentResolved { get; set; }
        public bool AdminActionRequired { get; set; }
    }

    public class Exam
    {
        private ExamResult _examResult;
        public Guid ExamId { get; set; }
        public DateTime ExamDeadline { get; set; }
        public int Score { get; set; }

        public ExamResult Result
        {
            get
            {
                if (Score >= 76)
                {
                    return _examResult = ExamResult.Pass;
                }

                if (Score >= 1 & Score < 76)
                {
                    return _examResult = ExamResult.Fail;
                }

                return _examResult = ExamResult.DidNotSitExam;
            }

            set => _examResult = value;
        }
        
        public bool NotificationSent { get; set; }

        public Exam()
        {
        }

        public Exam(Guid examId, DateTime examDeadline, int score)
        {
            ExamId = examId;
            ExamDeadline = examDeadline;
            Score = score;
            NotificationSent = false;
        }

        public override string ToString()
        {
            return
                $"ExamId: {ExamId}, ExamDate: {ExamDeadline}, Score: {Score}, Result: {Result}, NotificationSent: {NotificationSent}";
        }
    }

    public enum ExamResult
    {
        Pass,
        Fail,
        DidNotSitExam
    }
}