let refreshPromise: Promise<boolean> | null = null

async function refreshAuth(): Promise<boolean> {
  if (!refreshPromise) {
    refreshPromise = fetch('/auth/refresh', {
      method: 'POST',
      credentials: 'include',
    })
      .then((response) => response.ok)
      .finally(() => {
        refreshPromise = null
      })
  }

  return refreshPromise
}

export async function apiFetch(
  input: RequestInfo | URL,
  init: RequestInit = {},
): Promise<Response> {
  let response = await fetch(input, {
    ...init,
    credentials: 'include',
  })

  if (response.status !== 401) {
    return response
  }


  const refreshed = await refreshAuth()

  if (!refreshed) {
    window.location.href = '/auth/login'
    return response
  }

  response = await fetch(input, {
    ...init,
    credentials: 'include',
  })

  return response
}