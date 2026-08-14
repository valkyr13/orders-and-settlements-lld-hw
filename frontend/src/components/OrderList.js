import React, { useEffect, useState } from "react";
import { getOrders } from "../api";

function OrderList() {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        async function loadOrders() {
            try {
                const data = await getOrders();

                console.log("Orders:", data);

                setOrders(data.orders || []);
            } catch (err) {
                console.error(err);
                setError(err.message);
            } finally {
                setLoading(false);
            }
        }

        loadOrders();
    }, []);

    if (loading) {
        return React.createElement(
            "p",
            null,
            "Loading orders..."
        );
    }

    if (error) {
        return React.createElement(
            "p",
            null,
            `Error: ${error}`
        );
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
            "Orders"
        ),

        React.createElement(
            "p",
            null,
            `Found ${orders.length} orders`
        ),

        orders.map((order) =>
            React.createElement(
                "div",
                {
                    key: order.id,
                    style: {
                        padding: "10px",
                        borderBottom: "1px solid #ddd",
                    },
                },
                order.customer_name
            )
        )
    );
}

export default OrderList;