import type { ButtonHTMLAttributes, HTMLAttributes } from "vue"
import type { ButtonVariants } from "."
import { defineComponent } from "vue"
import { Primitive } from "reka-ui"
import { cn } from "@/lib/utils"
import { buttonVariants } from "."

interface ButtonProps {
  variant?: ButtonVariants["variant"]
  size?: ButtonVariants["size"]
  class?: HTMLAttributes["class"]
  as?: any
  asChild?: boolean
}

export default defineComponent<ButtonProps & Omit<ButtonHTMLAttributes, keyof ButtonProps>>({
  name: 'Button',
  props: ['variant', 'size', 'class', 'as', 'asChild', 'type', 'disabled'] as any,
  setup(props, { slots, attrs }) {
    return () => (
      <Primitive
        {...attrs}
        data-slot="button"
        as={props.as || "button"}
        asChild={props.asChild}
        class={cn(buttonVariants({ variant: props.variant, size: props.size }), props.class)}
      >
        {slots.default?.()}
      </Primitive>
    )
  },
})
