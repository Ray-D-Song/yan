import type { JsonArray, JsonObject, JsonValue } from './types'

/**
 * Custom fetcher options that extend native fetch to support Object body
 */
interface FetcherOptions extends Omit<RequestInit, 'body'> {
  body?: BodyInit | JsonObject | JsonArray
}

/**
 * Custom fetcher parameters that support Object body
 */
type FetcherParams = [input: RequestInfo | URL, init?: FetcherOptions]

/**
 * A robust fetch wrapper function.
 * This fetcher:
 * 1. Automatically injects the 'Authorization' header if a token exists in localStorage.
 * 2. Allows passing custom fetch options (method, body, custom headers, etc.).
 * 3. Automatically JSON.stringifies object bodies if Content-Type is 'application/json'.
 * 4. Automatically removes 'Content-Type' for FormData to let the browser handle it.
 * 5. Handles both HTTP errors (e.g., 404, 500) and API business logic errors (e.g., code !== 200).
 * 6. Automatically downloads binary files (images, PDFs) when detected.
 * @template T The expected type of the 'data' field in a successful API response.
 * @param args The same arguments as the native 'fetch' function (URL and RequestInit), with extended body support.
 * @returns {Promise<T>} A promise that resolves with the 'data' from the API response.
 * @throws {string} Throws an error message (from 'msg' or HTTP status) if the request fails.
 */
async function fetcher<T = JsonValue>(...args: FetcherParams): Promise<T> {
  // 1. Separate the arguments: 'input' is the URL, 'init' is the options object.
  const [input, init = {}] = args

  // 2. Retrieve the token from localStorage.
  const token = localStorage.getItem('token')

  // 3. Define default headers.
  const defaultHeaders: HeadersInit = {
    // We assume JSON APIs by default.
    'Organ-Code': localStorage.getItem('organ-code') || '',
    'Content-Type': 'application/json',
    // If a token exists, automatically add the Authorization header.
    ...(token && { Authorization: `${token}` }),
  }

  // 4. Create base options without body first
  const { body, ...initWithoutBody } = init

  // 5. Merge options:
  //    - User's 'init' options (like 'method', 'body') take precedence.
  //    - Headers need to be merged deeply.
  const finalOptions: RequestInit = {
    ...initWithoutBody, // User's options (e.g., method, credentials)
    headers: {
      ...defaultHeaders, // Our default headers
      ...(init.headers || {}), // User's custom headers (can override defaults)
    },
  }

  // 5. [Convenience Feature 1] Handle FormData automatically.
  //    If the body is FormData, the browser must set the 'Content-Type'
  //    (including the 'boundary' string) itself.
  //    We must delete our default 'Content-Type' header.
  if (body instanceof FormData) {
    delete (finalOptions.headers as Record<string, string>)['Content-Type']
    finalOptions.body = body
  }
  // 6. [Convenience Feature 2] Auto-serialize JSON body.
  //    If the body is a plain JS object (and not FormData, Blob, etc.)
  //    and the Content-Type is still 'application/json', serialize it.
  else if (body) {
    const contentType = (finalOptions.headers as Record<string, string>)['Content-Type']
    if (typeof body === 'object'
      && !(body instanceof Blob)
      && !(body instanceof URLSearchParams)
      && !(body instanceof ArrayBuffer)
      && contentType === 'application/json'
    ) {
      finalOptions.body = JSON.stringify(body)
    }
    else {
      finalOptions.body = body as BodyInit
    }
  }

  // 7. Make the request using the final URL and options.
  return fetch(`/api${input}`, finalOptions).then(async (res) => {
    // 8. [Improvement] Check for HTTP errors (non-2xx responses).
    if (!res.ok) {
      // Handle 401 Unauthorized - redirect to login page
      if (res.status === 401) {
        handleAuthError()
        throw new Error(await res.text())
      }

      // For non-2xx responses, read the response body as text (error message)
      const errorMessage = await res.text()
      throw new Error(errorMessage || res.statusText)
    }

    // Handle 204 No Content (successful, but no body to parse).
    if (res.status === 204) {
      return null // Or {} as T, depending on your API spec.
    }

    const contentType = res.headers.get('Content-Type')
    if (!contentType)
      return res.text()
    if (contentType === 'application/json') {
      return await res.json()
    }
    if (shouldDownloadAsFile(contentType)) {
      // Handle binary file download
      const blob = await res.blob()
      const filename = getFilenameFromHeaders(res.headers) || getFilenameFromUrl(input.toString()) || 'download'

      // Create download link
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = filename
      link.style.display = 'none'

      // Trigger download
      document.body.appendChild(link)
      link.click()

      // Cleanup
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)

      return blob as T
    }
  })
}

/**
 * Determine if content type should be downloaded as file
 * @param contentType MIME content type
 * @returns true if should download as file
 */
function shouldDownloadAsFile(contentType: string): boolean {
  const downloadableTypes = [
    'image/', // All image types
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'text/csv',
    'application/zip',
    'application/x-zip-compressed',
    'application/octet-stream',
  ]

  return downloadableTypes.some(type => contentType.includes(type))
}

export {
  fetcher,
}

/**
 * Extract filename from Content-Disposition header
 * @param headers Response headers
 * @returns Filename or null if not found
 */
function getFilenameFromHeaders(headers: Headers): string | null {
  const contentDisposition = headers.get('Content-Disposition')
  if (!contentDisposition)
    return null

  // Try to extract filename from Content-Disposition header
  // Common formats: attachment; filename="file.pdf" or attachment; filename=file.pdf
  const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((["']).*?\2|[^;\n]*)/)
  if (filenameMatch && filenameMatch[1]) {
    const filename = filenameMatch[1].replace(/['"]/g, '').trim()
    // Decode URL-encoded characters
    try {
      return decodeURIComponent(filename)
    }
    catch {
      // If decoding fails, return original filename
      return filename
    }
  }

  return null
}

/**
 * Extract filename from URL
 * @param url URL string
 * @returns Filename or null if not found
 */
function getFilenameFromUrl(url: string): string | null {
  try {
    const urlObj = new URL(url)
    const pathname = urlObj.pathname
    const filename = pathname.split('/').pop()
    return filename || null
  }
  catch {
    return null
  }
}

function handleAuthError() {
  // Clear all localStorage data
  localStorage.clear()
  // Redirect to login page only if not already on it
  if (window.location.pathname !== '/login') {
    window.location.href = '/login'
  }
}
