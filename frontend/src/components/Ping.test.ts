import { expect, test, vi } from 'vitest'
import { render } from 'vitest-browser-vue'
import Ping from './Ping.vue'

test('Ping component fetches and displays result', async () => {
  // Mock fetch
  const mockFetch = vi.spyOn(window, 'fetch').mockResolvedValue({
    ok: true,
    status: 200,
    text: () => Promise.resolve('pong'),
  } as Response)

  const { getByText } = render(Ping)

  // Should show loading initially
  await expect.element(getByText('Loading...')).toBeInTheDocument()

  // Wait for the result to appear
  await expect.element(getByText('Result: pong')).toBeInTheDocument()
  
  expect(mockFetch).toHaveBeenCalled()
  
  mockFetch.mockRestore()
})

test('Ping component handles error', async () => {
  // Mock fetch error
  const mockFetch = vi.spyOn(window, 'fetch').mockResolvedValue({
    ok: false,
    status: 500,
  } as Response)

  const { getByText } = render(Ping)

  // Should show error message
  await expect.element(getByText(/Failed to fetch ping/)).toBeInTheDocument()
  
  mockFetch.mockRestore()
})
