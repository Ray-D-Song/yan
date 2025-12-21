import AuthLayout from "@/layouts/auth-layout"
import { defineComponent } from "vue"
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

export default defineComponent(() => {
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
            <form>
              <FieldGroup>
                <Field>
                  <FieldLabel>
                    Email
                  </FieldLabel>
                  <Input
                    type="email"
                    placeholder="m@example.com"
                    required
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
                  <Input type="password" required />
                </Field>
                <Field>
                  <Button type="submit">
                    Login
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
