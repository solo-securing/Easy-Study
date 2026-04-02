"use client";

type ProgressTrackerProps = {
  completionPercent: number;
  completedUnitCount: number;
  officialQuizScore?: number;
};

export default function ProgressTracker({
  completionPercent,
  completedUnitCount,
  officialQuizScore,
}: ProgressTrackerProps) {
  const normalized = Math.max(0, Math.min(100, completionPercent));

  return (
    <section>
      <h2>Learning Progress</h2>
      <p>Completion: {normalized.toFixed(1)}%</p>
      <progress max={100} value={normalized} />
      <p>Completed units: {completedUnitCount}</p>
      <p>Official quiz score: {officialQuizScore ?? 0}</p>
    </section>
  );
}
