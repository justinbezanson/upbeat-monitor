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
  })
}))

describe('LoginView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  test('renders login form', async () => {
    const { getByLabelText, getByRole } = render(LoginView)

    await expect.element(getByLabelText('Email:')).toBeInTheDocument()
    await expect.element(getByLabelText('Password:')).toBeInTheDocument()
    await expect.element(getByRole('button', { name: 'Login' })).toBeInTheDocument()
  })

  test('calls authStore.login and redirects on success', async () => {
    const authStore = useAuthStore()
    // Mock login to succeed
    vi.spyOn(authStore, 'login').mockResolvedValue(true)

    const { getByLabelText, getByRole } = render(LoginView)

    const emailInput = getByLabelText('Email:')
    const passwordInput = getByLabelText('Password:')
    const submitButton = getByRole('button', { name: 'Login' })

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

    const { getByText, getByRole, getByLabelText } = render(LoginView)

    await getByLabelText('Email:').fill('test@example.com')
    await getByLabelText('Password:').fill('wrong')
    await getByRole('button', { name: 'Login' }).click()

    await expect.element(getByText('Invalid credentials')).toBeInTheDocument()
  })
})
