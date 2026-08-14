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

export function getCurrentUser() {
    return request("/auth/me");
}


export function getOrders() {
    return request("/orders");
}

export function getOrder(orderId) {
    return request(`/orders/${orderId}`);
}

export function createOrder(order) {
    return request("/orders", {
        method: "POST",
        body: JSON.stringify(order),
    });
}

export function updateOrder(orderId, order) {
    return request(`/orders/${orderId}`, {
        method: "PUT",
        body: JSON.stringify(order),
    });
}

export function deleteOrder(orderId) {
    return request(`/orders/${orderId}`, {
        method: "DELETE",
    });
}