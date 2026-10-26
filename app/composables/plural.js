// Czech plural forms: czPlural(1, "otázka", "otázky", "otázek") -> "otázka",
// 2–4 -> "otázky", 0 and 5+ -> "otázek". Decimals take the genitive singular
// ("1,5 bodu"), passed as `decimal`; it defaults to `few`, which matches the
// genitive singular for feminine nouns ("1,5 otázky").
export const czPlural = (count, one, few, many, decimal = few) => {
  const n = Math.abs(Number(count));
  if (!Number.isInteger(n)) return decimal;
  if (n === 1) return one;
  if (n >= 2 && n <= 4) return few;
  return many;
};

// "3 otázky": the count followed by the matching form.
export const czCount = (count, one, few, many, decimal = few) =>
  `${count} ${czPlural(count, one, few, many, decimal)}`;
