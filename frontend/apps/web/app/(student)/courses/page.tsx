import ProgressTracker from "../../../src/components/student/ProgressTracker";

type CourseSummary = {
  id: string;
  title: string;
  status: "draft" | "published" | "archived";
};

async function fetchMyCourses(): Promise<CourseSummary[]> {
  const response = await fetch("http://localhost:8080/api/v1/me/courses", {
    headers: {
      "X-Tenant-ID": "tenant-001",
    },
    cache: "no-store",
  });
  if (!response.ok) {
    return [];
  }
  const json = (await response.json()) as { items?: CourseSummary[] };
  return json.items ?? [];
}

export default async function StudentCoursesPage() {
  const courses = await fetchMyCourses();

  return (
    <main>
      <h1>My Courses</h1>
      {courses.length === 0 ? (
        <p>No assigned courses.</p>
      ) : (
        <ul>
          {courses.map((course) => (
            <li key={course.id}>
              <strong>{course.title}</strong> - {course.status}
            </li>
          ))}
        </ul>
      )}
      <ProgressTracker completionPercent={0} completedUnitCount={0} officialQuizScore={0} />
    </main>
  );
}
