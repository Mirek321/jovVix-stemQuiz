import { czPlural } from "~/composables/plural";

export function useGetTime(targetTime) {
  const quizCreationTime = new Date(targetTime);

  const currentTimeInMs = Date.now();

  const timeSinceCreation = currentTimeInMs - quizCreationTime.getTime();

  const secondsSinceCreation = Math.round(timeSinceCreation / 1000);
  const minutesSinceCreation = Math.round(secondsSinceCreation / 60);
  const hoursSinceCreation = Math.round(minutesSinceCreation / 60);
  const daysSinceCreation = Math.round(hoursSinceCreation / 24);
  const monthsSinceCreation = Math.round(daysSinceCreation / 30);
  const yearsSinceCreation = Math.round(monthsSinceCreation / 12);

  // "před" + instrumental: 1 → singular form, 2+ → plural form
  const ago = (n, one, other) =>
    `Vytvořeno před ${n} ${czPlural(n, one, other, other)}`;

  if (secondsSinceCreation < 60) {
    return "Vytvořeno právě teď";
  } else if (minutesSinceCreation < 60) {
    return ago(minutesSinceCreation, "minutou", "minutami");
  } else if (hoursSinceCreation < 24) {
    return ago(hoursSinceCreation, "hodinou", "hodinami");
  } else if (daysSinceCreation < 30) {
    return ago(daysSinceCreation, "dnem", "dny");
  } else if (monthsSinceCreation < 12) {
    return ago(monthsSinceCreation, "měsícem", "měsíci");
  }
  return ago(yearsSinceCreation, "rokem", "lety");
}
