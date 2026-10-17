import { usePush } from "notivue";
const toast = usePush();

// navigator.clipboard only exists on secure origins (https, localhost). At an
// event the app is often opened over plain http on a LAN address, so fall
// back to a hidden textarea + execCommand there.
const legacyCopy = (text) => {
  const area = document.createElement("textarea");
  area.value = text;
  area.setAttribute("readonly", "");
  area.style.position = "fixed";
  area.style.opacity = "0";
  document.body.appendChild(area);
  area.select();
  try {
    return document.execCommand("copy");
  } finally {
    document.body.removeChild(area);
  }
};

export default function usecopyToClipboard(text) {
  const copied = navigator.clipboard?.writeText
    ? navigator.clipboard.writeText(text)
    : Promise.resolve().then(() => {
        if (!legacyCopy(text)) throw new Error("execCommand copy failed");
      });

  copied
    .then(() => {
      toast.success("Zkopírováno do schránky");
    })
    .catch((error) => {
      toast.warning("Kopírování do schránky se nezdařilo", error);
    });
}
