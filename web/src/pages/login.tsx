import AuthLayout from "@/layouts/auth-layout"
import { defineComponent, reactive } from "vue"
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Field,
  FieldDescription,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { usersApi } from '@/api/users'
import { isDev } from '@/lib/env'
import { toast } from 'vue-sonner'

export default defineComponent(() => {
  const formData = reactive({
    email: isDev ? 'test@example.com' : '',
    password: isDev ? 'test123456' : '',
    isSubmitting: false,
  })

  const handleSubmit = async (e: Event) => {
    e.preventDefault()

    try {
      formData.isSubmitting = true
      const user = await usersApi.login({
        email: formData.email,
        password: formData.password,
      })

      localStorage.setItem('user', JSON.stringify(user))

      // Success - show toast and redirect after 1 second
      toast.success('Login successful! Redirecting...')
      setTimeout(() => {
        window.location.href = '/'
      }, 1000)
    } catch (error) {
      console.error(error)
      // Error - show toast
      const errorMessage = error instanceof Error ? error.message : String(error)
      toast.error(`Login failed: ${errorMessage}`)
    } finally {
      formData.isSubmitting = false
    }
  }

  return () => (
    <AuthLayout>
      <div class={cn('flex flex-col gap-6')}>
        <Card class="border-none shadow-none">
          <CardHeader>
            <CardTitle>Login to your account</CardTitle>
            <CardDescription>
              Enter your email below to login to your account
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit}>
              <FieldGroup>
                <Field>
                  <FieldLabel>
                    Email
                  </FieldLabel>
                  <Input
                    type="email"
                    placeholder="m@example.com"
                    required
                    v-model={formData.email}
                  />
                </Field>
                <Field>
                  <div class="flex items-center">
                    <FieldLabel>
                      Password
                    </FieldLabel>
                    <a
                      href="#"
                      class="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                    >
                      Forgot your password?
                    </a>
                  </div>
                  <Input
                    type="password"
                    required
                    v-model={formData.password}
                  />
                </Field>
                <Field>
                  <Button type="submit" disabled={formData.isSubmitting}>
                    {formData.isSubmitting ? 'Logging in...' : 'Login'}
                  </Button>
                  <FieldDescription class="text-center">
                    Don't have an account?
                    {' '}
                    <a href="/signup">
                      Sign up
                    </a>
                  </FieldDescription>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  )
})
