import AuthLayout from "@/layouts/auth-layout"
import { defineComponent, reactive, watch } from "vue"
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

// Generate random test data for development
function generateTestData() {
  const randomId = Math.random().toString(36).substring(2, 8)
  const testPassword = 'test123456'

  return {
    username: `testuser_${randomId}`,
    email: `test_${randomId}@example.com`,
    password: testPassword,
    confirmPassword: testPassword,
  }
}

export default defineComponent(() => {
  const formData = reactive({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    isSubmitting: false,
    ...(isDev ? generateTestData() : {}),
  })

  const handleSubmit = async (e: Event) => {
    e.preventDefault()

    // Validate password match
    if (formData.password !== formData.confirmPassword) {
      alert('Passwords do not match')
      return
    }

    // Validate password length
    if (formData.password.length < 6) {
      alert('Password must be at least 6 characters long')
      return
    }

    try {
      formData.isSubmitting = true
      await usersApi.register({
        username: formData.username,
        email: formData.email,
        password: formData.password,
      })

      // Success - show toast and redirect after 1 second
      toast.success('Registration successful! Redirecting to login...')
      setTimeout(() => {
        window.location.href = '/login'
      }, 1000)
    } catch (error) {
      // Error - show alert
      const errorMessage = error instanceof Error ? error.message : String(error)
      alert(`Registration failed: ${errorMessage}`)
    } finally {
      formData.isSubmitting = false
    }
  }

  return () => (
    <AuthLayout>
      <Card class="border-none shadow-none">
        <CardHeader>
          <CardTitle>Create an account</CardTitle>
          <CardDescription>
            Enter your information below to create your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <FieldGroup>
              <Field>
                <FieldLabel>
                  Full Name
                </FieldLabel>
                <Input
                  type="text"
                  placeholder="Ray Song"
                  required
                  v-model={formData.username}
                />
              </Field>
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
                <FieldDescription>
                  We'll use this to contact you. We will not share your email with
                  anyone else.
                </FieldDescription>
              </Field>
              <Field>
                <FieldLabel>
                  Password
                </FieldLabel>
                <Input
                  type="password"
                  required
                  v-model={formData.password}
                />
                <FieldDescription>Must be at least 6 characters long.</FieldDescription>
              </Field>
              <Field>
                <FieldLabel>
                  Confirm Password
                </FieldLabel>
                <Input
                  type="password"
                  required
                  v-model={formData.confirmPassword}
                />
                <FieldDescription>Please confirm your password.</FieldDescription>
              </Field>
              <FieldGroup>
                <Field>
                  <Button type="submit" disabled={formData.isSubmitting}>
                    {formData.isSubmitting ? 'Creating Account...' : 'Create Account'}
                  </Button>
                  <FieldDescription class="px-6 text-center">
                    Already have an account? <a href="/login">Sign in</a>
                  </FieldDescription>
                </Field>
              </FieldGroup>
            </FieldGroup>
          </form>
        </CardContent>
      </Card>
    </AuthLayout>
  )
})
