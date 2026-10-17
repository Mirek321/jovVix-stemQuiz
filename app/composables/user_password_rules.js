import { computed } from "vue";

export function useUserPasswordRules(passwordRef, firstnameRef, lastnameRef) {
  const rules = {
    minLength: 8,
    alphanumeric: /^(?=.*[a-zA-Z])(?=.*\d)/,
    specialChar: /[^a-zA-Z0-9]/,
  };

  const passwordErrors = computed(() => {
    const errors = [];

    if (!passwordRef.value) {
      return errors;
    }

    if (passwordRef.value.length < rules.minLength) {
      errors.push("Heslo musí mít aspoň 8 znaků.");
    }

    if (!rules.alphanumeric.test(passwordRef.value)) {
      errors.push("Heslo musí obsahovat písmena i číslice.");
    }

    if (!rules.specialChar.test(passwordRef.value)) {
      errors.push("Heslo musí obsahovat aspoň jeden speciální znak.");
    }

    const first = firstnameRef?.value?.toLowerCase();
    const last = lastnameRef?.value?.toLowerCase();

    if (
      (first && passwordRef.value.toLowerCase().includes(first)) ||
      (last && passwordRef.value.toLowerCase().includes(last))
    ) {
      errors.push("Heslo nesmí obsahovat tvoje jméno.");
    }

    return errors;
  });

  return {
    passwordErrors,
  };
}
