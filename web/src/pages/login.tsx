import LoginForm from "@/components/LoginForm.vue";
import AuthLayout from "@/layouts/auth-layout";
import { defineComponent } from "vue";

export default defineComponent(() => {
  return () => (
    <AuthLayout>
      <LoginForm />
    </AuthLayout>
  )
})
