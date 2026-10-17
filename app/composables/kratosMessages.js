// Czech texts for Ory Kratos UI messages, keyed by Kratos message id.
// Kratos always answers in English; unknown ids fall back to message.text.
// See https://www.ory.sh/docs/kratos/concepts/ui-messages
const KRATOS_MESSAGES_CS = {
  // Info – login / registration
  1010001: "Přihlásit se",
  1010003: "Potvrď prosím, že jsi to opravdu ty.",
  1040001: "Zaregistrovat se",
  // Info – settings
  1050001: "Změny byly uloženy.",
  // Info – recovery
  1060001: "Účet se podařilo obnovit. Teď si prosím nastav nové heslo.",
  1060002:
    "Na zadaný e-mail jsme poslali odkaz pro obnovení účtu. Pokud nic nepřišlo, zkontroluj, že je adresa správně a že ji používáš u svého účtu.",
  1060003:
    "Na zadaný e-mail jsme poslali kód pro obnovení účtu. Pokud nic nepřišlo, zkontroluj, že je adresa správně a že ji používáš u svého účtu.",
  // Info – verification
  1080001: "Na tvůj e-mail jsme poslali ověřovací odkaz.",
  1080002: "E-mail se podařilo ověřit.",
  1080003: "Na tvůj e-mail jsme poslali ověřovací kód.",
  // Validation errors
  4000002: "Toto pole je povinné.",
  4000005:
    "Toto heslo nelze použít. Zvol silnější heslo, které se nepodobá tvému e-mailu.",
  4000006: "Zadaný e-mail nebo heslo není správné.",
  4000007: "Účet s tímto e-mailem už existuje.",
  4000010: "Účet ještě není ověřený. Nejdřív si prosím ověř e-mail.",
  // Expired flows
  4010001: "Přihlášení vypršelo. Zkus to prosím znovu.",
  4040001: "Registrace vypršela. Zkus to prosím znovu.",
  4050001: "Platnost formuláře vypršela. Zkus to prosím znovu.",
  4060005: "Obnovení účtu vypršelo. Požádej prosím o nový kód.",
  4070005: "Ověření vypršelo. Požádej prosím o nový kód.",
  // Invalid / used codes
  4060004: "Odkaz pro obnovení je neplatný nebo už byl použit.",
  4060006: "Kód pro obnovení je neplatný nebo už byl použit.",
  4070001: "Ověřovací odkaz je neplatný nebo už byl použit.",
  4070006: "Ověřovací kód je neplatný nebo už byl použit.",
};

export const kratosMessageText = (message) => {
  if (!message) return "";
  return KRATOS_MESSAGES_CS[message.id] || message.text || "";
};

// Joins all messages of a Kratos UI node into one Czech string.
export const kratosNodeMessagesText = (messages = []) =>
  messages.map(kratosMessageText).filter(Boolean).join(", ");
