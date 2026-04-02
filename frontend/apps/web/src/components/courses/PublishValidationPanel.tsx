"use client";

import { useMemo } from "react";

type ValidationIssue = {
  path: string;
  reason: string;
};

type PublishValidationPanelProps = {
  issues: ValidationIssue[];
};

export default function PublishValidationPanel({ issues }: PublishValidationPanelProps) {
  const hasIssues = useMemo(() => issues.length > 0, [issues]);

  if (!hasIssues) {
    return (
      <section>
        <h2>Publish Validation</h2>
        <p>Course structure is valid and ready to publish.</p>
      </section>
    );
  }

  return (
    <section>
      <h2>Publish Validation</h2>
      <p>Resolve these issues before publishing:</p>
      <ul>
        {issues.map((issue) => (
          <li key={`${issue.path}-${issue.reason}`}>
            <strong>{issue.path}</strong>: {issue.reason}
          </li>
        ))}
      </ul>
    </section>
  );
}
