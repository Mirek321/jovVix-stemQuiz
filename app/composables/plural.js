// Czech plural forms: czPlural(1, "otázka", "otázky", "otázek") -> "otázka",
// 2–4 -> "otázky", 0 and 5+ -> "otázek". Decimals take the "few" form.
export const czPlural = (count, one, few, many) => {
  const n = Math.abs(Number(count));
  if (!Number.isInteger(n)) return few;
  if (n === 1) return one;
  if (n >= 2 && n <= 4) return few;
  return many;
};

// "3 otázky": the count followed by the matching form.
export const czCount = (count, one, few, many) =>
  `${count} ${czPlural(count, one, few, many)}`;
