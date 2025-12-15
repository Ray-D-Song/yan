import { defineComponent } from "vue";

export default defineComponent((_, { slots }) => {
  return () => <div class="flex justify-center items-center h-full">
  <div class="max-w-sm w-full">
    {slots.default?.()}
  </div>
  </div>
})
