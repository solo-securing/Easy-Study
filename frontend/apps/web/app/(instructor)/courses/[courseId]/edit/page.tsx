"use client";

import { FormEvent, useMemo, useState } from "react";

import PublishValidationPanel from "../../../../../src/components/courses/PublishValidationPanel";

type Props = {
  params: {
    courseId: string;
  };
};

type SectionInput = {
  title: string;
  sortOrder: number;
  subsections: Array<{
    title: string;
    sortOrder: number;
    units: Array<{
      title: string;
      type: "video" | "document" | "quiz";
      sortOrder: number;
      contentRef?: string;
      quizId?: string;
    }>;
  }>;
};

export default function InstructorCourseEditorPage({ params }: Props) {
  const [tenantId, setTenantId] = useState("tenant-001");
  const [title, setTitle] = useState("Course title");
  const [description, setDescription] = useState("");
  const [statusMessage, setStatusMessage] = useState("");
  const [issues, setIssues] = useState<Array<{ path: string; reason: string }>>([]);

  const defaultStructure = useMemo<SectionInput[]>(
    () => [
      {
        title: "Section 1",
        sortOrder: 0,
        subsections: [
          {
            title: "Subsection 1",
            sortOrder: 0,
            units: [
              {
                title: "Unit 1",
                type: "video",
                sortOrder: 0,
                contentRef: "asset://video-1",
              },
            ],
          },
        ],
      },
    ],
    [],
  );

  const saveCourse = async (event: FormEvent) => {
    event.preventDefault();
    setStatusMessage("");
    setIssues([]);

    const patchRes = await fetch(`http://localhost:8080/api/v1/tenants/${tenantId}/courses/${params.courseId}`, {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": tenantId,
      },
      body: JSON.stringify({ title, description }),
    });
    if (!patchRes.ok) {
      setStatusMessage("Failed to update course metadata");
      return;
    }

    const structureRes = await fetch(
      `http://localhost:8080/api/v1/tenants/${tenantId}/courses/${params.courseId}/structure`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          "X-Tenant-ID": tenantId,
        },
        body: JSON.stringify({ sections: defaultStructure }),
      },
    );
    if (!structureRes.ok) {
      setStatusMessage("Failed to update course structure");
      return;
    }

    setStatusMessage("Course metadata and structure saved.");
  };

  const publishCourse = async () => {
    setStatusMessage("");
    setIssues([]);

    const response = await fetch(
      `http://localhost:8080/api/v1/tenants/${tenantId}/courses/${params.courseId}/publish`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Tenant-ID": tenantId,
        },
      },
    );

    if (response.status === 422) {
      const payload = (await response.json()) as { missingNodes?: Array<{ path: string; reason: string }> };
      setIssues(payload.missingNodes ?? []);
      setStatusMessage("Publish blocked by validation issues.");
      return;
    }
    if (!response.ok) {
      setStatusMessage("Failed to publish course.");
      return;
    }
    setStatusMessage("Course published.");
  };

  return (
    <main>
      <h1>Course Editor</h1>
      <p>Course ID: {params.courseId}</p>
      <form onSubmit={saveCourse}>
        <label>
          Tenant ID
          <input value={tenantId} onChange={(event) => setTenantId(event.target.value)} />
        </label>
        <label>
          Title
          <input value={title} onChange={(event) => setTitle(event.target.value)} />
        </label>
        <label>
          Description
          <textarea value={description} onChange={(event) => setDescription(event.target.value)} />
        </label>
        <button type="submit">Save Course</button>
      </form>
      <button onClick={publishCourse} type="button">
        Publish Course
      </button>
      <PublishValidationPanel issues={issues} />
      {statusMessage ? <p>{statusMessage}</p> : null}
    </main>
  );
}
