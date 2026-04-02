import { test, expect } from "@playwright/test";

test("student learning flow: my courses and quiz attempt page render", async ({ page }) => {
  await page.goto("/(student)/courses");
  await expect(page.getByRole("heading", { name: "My Courses" })).toBeVisible();

  await page.goto("/(student)/quizzes/quiz-001/attempt");
  await expect(page.getByRole("heading", { name: "Quiz Attempt" })).toBeVisible();
  await expect(page.getByRole("button", { name: "Start Attempt" })).toBeVisible();
});
