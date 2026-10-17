const fallbackImages = [
  "/images/landing/homepage-public-quiz-1.png",
  "/images/landing/homepage-public-quiz-2.png",
  "/images/landing/homepage-public-quiz-3.png",
  "/images/landing/homepage-public-quiz-4.png",
];

const hashIndex = (value, length) => {
  if (!value) return 0;
  const str = String(value);
  let hash = 0;
  for (let i = 0; i < str.length; i += 1) {
    hash = (hash * 31 + str.charCodeAt(i)) >>> 0;
  }
  return hash % length;
};

export const useQuizCover = () => {
  const { apiUrl } = useRuntimeConfig().public;

  const fallbackUrl = (quiz) =>
    fallbackImages[hashIndex(quiz?.id, fallbackImages.length)];

  const coverUrl = (quiz) =>
    quiz?.has_cover_image
      ? `${apiUrl}/quizzes/${quiz.id}/cover?v=${Math.floor(
          Date.parse(quiz.updated_at) / 1000
        )}`
      : fallbackUrl(quiz);

  return { coverUrl, fallbackUrl };
};
