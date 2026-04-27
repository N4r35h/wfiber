type HttpResponse<T = any> = {
	data: T
}

const baseURL = "http://localhost:8080/api"

const request = async (path: string, init?: RequestInit): Promise<HttpResponse> => {
	const response = await fetch(baseURL + path, {
		headers: {
			"Content-Type": "application/json",
			Accept: "application/json",
		},
		...init,
	})
	if (!response.ok) {
		throw new Error(`Request failed: ${response.status}`)
	}
	return {
		data: await response.json(),
	}
}

const http = {
	get: (path: string) => request(path),
	post: (path: string, body: unknown) =>
		request(path, {
			method: "POST",
			body: JSON.stringify(body),
		}),
}

export default http
