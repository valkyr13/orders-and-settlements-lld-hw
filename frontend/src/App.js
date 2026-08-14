import React, { useEffect, useState } from "react";
import Login from "./components/Login";
import { getCurrentUser } from "./api";

function App() {
    const [authenticated, setAuthenticated] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function checkSession() {
            try {
                await getCurrentUser();
                setAuthenticated(true);
            } catch {
                setAuthenticated(false);
            } finally {
                setLoading(false);
            }
        }

        checkSession();
    }, []);

    if (loading) {
        return React.createElement(
            "p",
            null,
            "Loading..."
        );
    }

    if (!authenticated) {
        return React.createElement(Login, {
            onLogin: () => setAuthenticated(true),
        });
    }

    return React.createElement(
        "div",
        { style: { padding: "40px" } },
        React.createElement("h1", null, "Welcome"),
        React.createElement(
            "p",
            null,
            "You are authenticated."
        )
    );
}

export default App;