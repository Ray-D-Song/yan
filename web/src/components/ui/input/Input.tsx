import type { HTMLAttributes, InputHTMLAttributes } from "vue"
import { defineComponent, reactive, watch } from "vue"
import { cn } from "@/lib/utils"

interface InputProps {
  defaultValue?: string | number
  modelValue?: string | number
  class?: HTMLAttributes["class"]
}

export default defineComponent<InputProps & Omit<InputHTMLAttributes, 'class' | 'defaultValue' | 'modelValue'>>({
  props: ['defaultValue', 'modelValue', 'class', 'type', 'placeholder', 'required', 'disabled', 'value'],
  emits: ['update:modelValue'],
  setup(props, { emit, attrs }) {
    const state = reactive({
      internalValue: props.modelValue ?? props.defaultValue ?? '',
    })

    watch(() => props.modelValue, (newVal) => {
      if (newVal !== undefined) {
        state.internalValue = newVal
      }
    })

    const handleInput = (e: Event) => {
      const target = e.target as HTMLInputElement
      state.internalValue = target.value
      emit('update:modelValue', target.value)
    }

    return () => (
      <input
        type={props.type}
        placeholder={props.placeholder}
        required={props.required}
        disabled={props.disabled}
        {...attrs}
        data-slot="input"
        class={cn(
          'file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground dark:bg-input/30 border-input h-9 w-full min-w-0 rounded-md border bg-transparent px-3 py-1 text-base shadow-xs transition-[color,box-shadow] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
          'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
          'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
          props.class,
        )}
        value={state.internalValue}
        onInput={handleInput}
      />
    )
  },
})
