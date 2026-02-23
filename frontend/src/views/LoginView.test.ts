import { expect, test, vi, beforeEach, describe } from 'vitest'
import { render } from 'vitest-browser-vue'
import { createPinia, setActivePinia } from 'pinia'
import LoginView from './LoginView.vue'
import { useAuthStore } from '@/stores/auth'

// Mock vue-router
const push = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push
  }),
  RouterLink: {
    template: '<a><slot /></a>'
  }
}))

describe('LoginView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  test('renders login form', async () => {
    const { getByLabelText, getByRole } = render(LoginView, {
        global: {
            stubs: {
                'router-link': true
            }
        }
    })

    await expect.element(getByLabelText('Email address')).toBeInTheDocument()
    await expect.element(getByLabelText('Password')).toBeInTheDocument()
    await expect.element(getByRole('button', { name: 'Sign in' })).toBeInTheDocument()
  })

  test('calls authStore.login and redirects on success', async () => {
    const authStore = useAuthStore()
    // Mock login to succeed
    vi.spyOn(authStore, 'login').mockResolvedValue(true)

    const { getByLabelText, getByRole } = render(LoginView, {
        global: {
            stubs: {
                'router-link': true
            }
        }
    })

    const emailInput = getByLabelText('Email address')
    const passwordInput = getByLabelText('Password')
    const submitButton = getByRole('button', { name: 'Sign in' })

    await emailInput.fill('test@example.com')
    await passwordInput.fill('password123')
    await submitButton.click()

    expect(authStore.login).toHaveBeenCalledWith('test@example.com', 'password123')
    expect(push).toHaveBeenCalledWith('/')
  })

  test('shows error message on failure', async () => {
    const authStore = useAuthStore()
    vi.spyOn(authStore, 'login').mockResolvedValue(false)
    authStore.error = 'Invalid credentials'

    const { getByText, getByRole, getByLabelText } = render(LoginView, {
        global: {
            stubs: {
                'router-link': true
            }
        }
    })

    await getByLabelText('Email address').fill('test@example.com')
    await getByLabelText('Password').fill('wrong')
    await getByRole('button', { name: 'Sign in' }).click()

    await expect.element(getByText('Invalid credentials')).toBeInTheDocument()
  })
})
