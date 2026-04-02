type CourseReportItem = {
  userId: string;
  completionPercent: number;
  completedUnitCount: number;
  officialQuizScore?: number;
};

type CourseReportResponse = {
  courseId: string;
  learnerCount: number;
  completionRate: number;
  items: CourseReportItem[];
};

async function fetchCourseReport(courseId: string): Promise<CourseReportResponse | null> {
  const response = await fetch(`http://localhost:8080/api/v1/tenants/tenant-001/reports/courses/${courseId}`, {
    headers: {
      "X-Tenant-ID": "tenant-001",
    },
    cache: "no-store",
  });
  if (!response.ok) {
    return null;
  }
  return (await response.json()) as CourseReportResponse;
}

export default async function InstructorCourseReportPage({ params }: { params: { courseId: string } }) {
  const report = await fetchCourseReport(params.courseId);

  return (
    <main>
      <h1>Course Report</h1>
      {!report ? (
        <p>No report data.</p>
      ) : (
        <>
          <p>Course ID: {report.courseId}</p>
          <p>Learners: {report.learnerCount}</p>
          <p>Completion rate: {report.completionRate}%</p>
          <ul>
            {report.items.map((item) => (
              <li key={item.userId}>
                {item.userId}: {item.completionPercent}% ({item.completedUnitCount} units)
              </li>
            ))}
          </ul>
        </>
      )}
    </main>
  );
}
