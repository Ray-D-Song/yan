import { defineComponent } from "vue";
import { RouterView } from "vue-router";
import 'vue-sonner/style.css'
import { Toaster } from '@/components/ui/sonner'

export default defineComponent(() => {
  return () => <>
    <RouterView></RouterView>
    <Toaster position="top-center" />
  </>
})
