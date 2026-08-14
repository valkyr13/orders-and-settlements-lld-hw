import React, { useState } from "react";
import Login from "./components/Login";

function App() {
    const [authenticated, setAuthenticated] = useState(false);

    if (!authenticated) {
        return React.createElement(Login, {
            onLogin: () => {
                setAuthenticated(true);
            },
        });
    }

    return React.createElement(
        "div",
        {
            style: {
                padding: "40px",
            },
        },
        React.createElement(
            "h1",
            null,
            "Welcome"
        ),
        React.createElement(
            "p",
            null,
            "You are authenticated."
        )
    );
}

export default App;