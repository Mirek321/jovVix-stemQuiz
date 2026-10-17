import { usePush } from "notivue";
const toast = usePush();

export default function usecopyToClipboard(text) {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      toast.success("Zkopírováno do schránky");
    })
    .catch((error) => {
      toast.warning("Kopírování do schránky se nezdařilo", error);
    });
}
