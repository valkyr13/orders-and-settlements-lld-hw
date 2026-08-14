import React, { useEffect, useState } from "react";
import { getOrder } from "../api";

function formatMoney(cents) {
    return `₹${(cents / 100).toFixed(2)}`;
}

function OrderDetails({ orderId, onBack }) {
    const [order, setOrder] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        async function loadOrder() {
            try {
                const data = await getOrder(orderId);
                setOrder(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        }

        loadOrder();
    }, [orderId]);

    if (loading) {
        return React.createElement("p", null, "Loading order...");
    }

    if (error) {
        return React.createElement(
            "div",
            null,
            React.createElement("p", null, `Error: ${error}`),
            React.createElement(
                "button",
                { onClick: onBack },
                "Back"
            )
        );
    }

    return React.createElement(
        "div",
        { style: { padding: "40px" } },

        React.createElement(
            "button",
            { onClick: onBack },
            "← Back"
        ),

        React.createElement(
            "h1",
            null,
            order.customer_name
        ),

        React.createElement(
            "p",
            null,
            `Due: ${new Date(order.due_date).toLocaleDateString()}`
        ),

        React.createElement(
            "p",
            null,
            `Status: ${order.status}`
        ),

        React.createElement(
            "h2",
            null,
            "Line Items"
        ),

        order.line_items.map((item) =>
            React.createElement(
                "div",
                {
                    key: item.id,
                    style: {
                        padding: "8px 0",
                    },
                },

                React.createElement(
                    "span",
                    null,
                    `${item.description} × ${item.quantity}`
                ),

                React.createElement(
                    "span",
                    {
                        style: {
                            marginLeft: "30px",
                        },
                    },
                    formatMoney(
                        item.quantity *
                        item.unit_price_cents
                    )
                )
            )
        ),

        React.createElement(
            "hr",
            null
        ),

        React.createElement(
            "p",
            null,
            `Total: ${formatMoney(order.total_cents)}`
        ),

        React.createElement(
            "p",
            null,
            `Paid: ${formatMoney(order.amount_paid_cents)}`
        ),

        React.createElement(
            "p",
            null,
            `Amount Due: ${formatMoney(order.amount_due_cents)}`
        )
    );
}

export default OrderDetails;