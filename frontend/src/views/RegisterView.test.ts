import { expect, test, vi, beforeEach, describe } from 'vitest'
import { render } from 'vitest-browser-vue'
import { createPinia, setActivePinia } from 'pinia'
import RegisterView from './RegisterView.vue'
import { useAuthStore } from '@/stores/auth'

// Mock vue-router
const push = vi.fn()
vi.mock('vue-router', () => ({
  useRouter: () => ({
    push
  })
}))

describe('RegisterView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  test('renders register form', async () => {
    const { getByLabelText, getByRole } = render(RegisterView)

    await expect.element(getByLabelText('Email:')).toBeInTheDocument()
    await expect.element(getByLabelText('Password:')).toBeInTheDocument()
    await expect.element(getByRole('button', { name: 'Register' })).toBeInTheDocument()
  })

  test('calls authStore.register and redirects on success', async () => {
    const authStore = useAuthStore()
    vi.spyOn(authStore, 'register').mockResolvedValue(true)

    const { getByLabelText, getByRole } = render(RegisterView)

    await getByLabelText('Email:').fill('newuser@example.com')
    await getByLabelText('Password:').fill('password123')
    await getByRole('button', { name: 'Register' }).click()

    expect(authStore.register).toHaveBeenCalledWith('newuser@example.com', 'password123')
    expect(push).toHaveBeenCalledWith('/login')
  })

  test('shows success message on successful registration', async () => {
    const authStore = useAuthStore()
    vi.spyOn(authStore, 'register').mockImplementation(async () => {
        authStore.successMessage = 'Registration successful!'
        return true
    })

    const { getByText, getByRole, getByLabelText } = render(RegisterView)

    await getByLabelText('Email:').fill('newuser@example.com')
    await getByLabelText('Password:').fill('password123')
    await getByRole('button', { name: 'Register' }).click()

    await expect.element(getByText('Registration successful!')).toBeInTheDocument()
  })

  test('shows error message on failure', async () => {
    const authStore = useAuthStore()
    vi.spyOn(authStore, 'register').mockResolvedValue(false)
    authStore.error = 'User already exists'

    const { getByText, getByRole, getByLabelText } = render(RegisterView)

    await getByLabelText('Email:').fill('existing@example.com')
    await getByLabelText('Password:').fill('password123')
    await getByRole('button', { name: 'Register' }).click()

    await expect.element(getByText('User already exists')).toBeInTheDocument()
  })
})
