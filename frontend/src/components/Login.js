    import React, { useState } from "react";
import { signup, login } from "../api";


function Login({ onLogin })  {
    const [isSignup, setIsSignup] = useState(false);

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    async function handleSubmit(event) {
    event.preventDefault();

    try {
        let response;

        if (isSignup) {
            response = await signup(email, password);
        }

        response = await login(email, password);

        console.log(
            "Authentication successful:",
            response
        );

        onLogin();
    } catch (error) {
        console.error(
            "Authentication failed:",
            error
        );
    }
}

    return React.createElement(
        "div",
        {
            style: {
                width: "400px",
                margin: "80px auto",
                padding: "24px",
            },
        },

        React.createElement(
            "h1",
            null,
            "Orders & Settlements"
        ),

        React.createElement(
            "h2",
            null,
            isSignup ? "Create Account" : "Login"
        ),

        React.createElement(
            "form",
            {
                onSubmit: handleSubmit,
            },

            React.createElement(
                "div",
                { style: { marginBottom: "16px" } },

                React.createElement(
                    "label",
                    null,
                    "Email"
                ),

                React.createElement("input", {
                    type: "email",
                    value: email,
                    onChange: (event) => {
                        setEmail(event.target.value);
                    },
                    required: true,
                    style: {
                        display: "block",
                        width: "100%",
                        padding: "8px",
                        marginTop: "4px",
                    },
                })
            ),

            React.createElement(
                "div",
                { style: { marginBottom: "16px" } },

                React.createElement(
                    "label",
                    null,
                    "Password"
                ),

                React.createElement("input", {
                    type: "password",
                    value: password,
                    onChange: (event) => {
                        setPassword(event.target.value);
                    },
                    required: true,
                    style: {
                        display: "block",
                        width: "100%",
                        padding: "8px",
                        marginTop: "4px",
                    },
                })
            ),

            React.createElement(
                "button",
                {
                    type: "submit",
                },
                isSignup ? "Sign Up" : "Login"
            )
        ),

        React.createElement(
            "button",
            {
                type: "button",
                onClick: () => {
                    setIsSignup(!isSignup);
                },
                style: {
                    marginTop: "16px",
                },
            },
            isSignup
                ? "Already have an account? Login"
                : "Create a new account"
        )
    );
}



export default Login;