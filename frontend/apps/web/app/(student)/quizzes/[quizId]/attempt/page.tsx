"use client";

import { useEffect, useMemo, useState } from "react";

type Props = {
  params: {
    quizId: string;
  };
};

export default function StudentQuizAttemptPage({ params }: Props) {
  const [tenantId, setTenantId] = useState("tenant-001");
  const [attemptId, setAttemptId] = useState("");
  const [answer, setAnswer] = useState("");
  const [status, setStatus] = useState("");
  const [timeLeft, setTimeLeft] = useState(1800);

  const autosavePayload = useMemo(
    () => ({
      answers: [{ questionId: "question-001", textAnswer: answer }],
      clientSavedAt: new Date().toISOString(),
    }),
    [answer],
  );

  useEffect(() => {
    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev <= 1) {
          clearInterval(timer);
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    if (!attemptId) return;
    const autosave = setInterval(async () => {
      await fetch(
        `http://localhost:8080/api/v1/tenants/${tenantId}/quizzes/${params.quizId}/attempts/${attemptId}/autosave`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
            "X-Tenant-ID": tenantId,
          },
          body: JSON.stringify(autosavePayload),
        },
      );
      setStatus("Autosaved");
    }, 10000);
    return () => clearInterval(autosave);
  }, [attemptId, autosavePayload, params.quizId, tenantId]);

  const startAttempt = async () => {
    const response = await fetch(`http://localhost:8080/api/v1/tenants/${tenantId}/quizzes/${params.quizId}/attempts`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": tenantId,
      },
      body: JSON.stringify({}),
    });
    if (!response.ok) {
      setStatus("Failed to start quiz attempt");
      return;
    }
    const json = (await response.json()) as { attemptId?: string };
    setAttemptId(json.attemptId ?? "");
    setStatus("Attempt started");
  };

  const submitAttempt = async () => {
    if (!attemptId) return;
    const response = await fetch(
      `http://localhost:8080/api/v1/tenants/${tenantId}/quizzes/${params.quizId}/attempts/${attemptId}/submit`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Tenant-ID": tenantId,
        },
      },
    );
    if (!response.ok) {
      setStatus("Submit failed");
      return;
    }
    setStatus("Submitted");
  };

  useEffect(() => {
    if (timeLeft > 0 || !attemptId) {
      return;
    }
    void submitAttempt();
  }, [timeLeft, attemptId]);

  return (
    <main>
      <h1>Quiz Attempt</h1>
      <p>Quiz: {params.quizId}</p>
      <label>
        Tenant ID
        <input value={tenantId} onChange={(event) => setTenantId(event.target.value)} />
      </label>
      <p>Time left: {timeLeft}s</p>
      <button onClick={startAttempt} type="button">
        Start Attempt
      </button>
      <textarea value={answer} onChange={(event) => setAnswer(event.target.value)} placeholder="Type your answer" />
      <button onClick={submitAttempt} type="button" disabled={!attemptId}>
        Submit
      </button>
      {status ? <p>{status}</p> : null}
    </main>
  );
}
