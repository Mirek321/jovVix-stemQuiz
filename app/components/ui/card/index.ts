import type { VariantProps } from "class-variance-authority";
import { cva } from "class-variance-authority";

export { default as Card } from "./Card.vue";
export { default as CardHeader } from "./CardHeader.vue";
export { default as CardBody } from "./CardBody.vue";
export { default as CardFooter } from "./CardFooter.vue";

export const cardVariants = cva("relative bg-jv-white text-jv-ink", {
  variants: {
    variant: {
      default: "jv-card border border-jv-line",
      rough: "jv-border-rough",
      even: "jv-border-even",
      uneven: "jv-border-uneven",
      flat: "border border-jv-line rounded-md",
      ghost: "",
    },
    tone: {
      canvas: "bg-jv-canvas",
      white: "bg-jv-white",
      yellow: "bg-jv-yellow",
      yellowSoft: "bg-jv-yellow-soft",
      mint: "bg-jv-mint",
      salmon: "bg-jv-salmon",
      lavender: "bg-jv-lavender",
      coral: "bg-jv-coral text-jv-white",
      ivory: "bg-jv-ivory",
    },
    shadow: {
      none: "",
      sm: "shadow-soft-sm",
      md: "shadow-soft",
      lg: "shadow-soft-lg",
    },
    padding: {
      none: "",
      sm: "p-3 sm:p-4",
      md: "p-4 sm:p-6",
      lg: "p-6 sm:p-8",
    },
  },
  defaultVariants: {
    variant: "default",
    tone: "white",
    shadow: "md",
    padding: "md",
  },
});

export type CardVariants = VariantProps<typeof cardVariants>;
