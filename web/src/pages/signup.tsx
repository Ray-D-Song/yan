import SignupForm from "@/components/SignupForm.vue";
import { defineComponent } from "vue";
import AuthLayout from "@/layouts/auth-layout";

export default defineComponent(() => {
  return () => (
    <AuthLayout>
      <SignupForm />
    </AuthLayout>
  )
})
