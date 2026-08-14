const API_BASE_URL = "http://localhost:8080";

async function request(endpoint, options = {}) {
    const response = await fetch(
        `${API_BASE_URL}${endpoint}`,
        {
            credentials: "include",

            headers: {
                "Content-Type": "application/json",
                ...(options.headers || {}),
            },

            ...options,
        }
    );

    const data = await response.json().catch(() => ({}));

    if (!response.ok) {
        throw new Error(
            data.error || "Request failed"
        );
    }

    return data;
}

export function signup(email, password) {
    return request("/auth/signup", {
        method: "POST",

        body: JSON.stringify({
            email,
            password,
        }),
    });
}

export function login(email, password) {
    return request("/auth/login", {
        method: "POST",

        body: JSON.stringify({
            email,
            password,
        }),
    });
}