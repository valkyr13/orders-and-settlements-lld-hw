import React, { useEffect, useState } from "react";
import { getOrders } from "../api";

function formatMoney(cents) {
    return `₹${(cents / 100).toFixed(2)}`;
}

function formatDate(date) {
    return new Date(date).toLocaleDateString();
}

function OrderList({
    onSelectOrder,
    onCreateOrder,
}) {    
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        async function loadOrders() {
            try {
                const data = await getOrders();
                setOrders(data.orders || []);
            } catch (err) {
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
        { style: { padding: "40px" } },

        React.createElement(
            "div",
            {
                style: {
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                    marginBottom: "24px",
                },
            },

            React.createElement(
                "div",
                null,
                React.createElement(
                    "h1",
                    null,
                    "Orders"
                ),

                React.createElement(
                    "p",
                    null,
                    `${orders.length} order${orders.length === 1 ? "" : "s"}`
                ),

                React.createElement( "button",{
                        onClick: onCreateOrder,},"+ Create Order"))),

        orders.length === 0
            ? React.createElement(
                  "p",
                  null,
                  "No orders found."
              )
            : React.createElement(
                  "table",
                  {
                      style: {
                          width: "100%",
                          borderCollapse: "collapse",
                      },
                  },

                  React.createElement(
                      "thead",
                      null,

                      React.createElement(
                          "tr",
                          null,
                          React.createElement("th", null, "Customer"),
                          React.createElement("th", null, "Due Date"),
                          React.createElement("th", null, "Total"),
                          React.createElement("th", null, "Paid"),
                          React.createElement("th", null, "Due"),
                          React.createElement("th", null, "Status")
                      )
                  ),

                  React.createElement(
                      "tbody",
                      null,

                      orders.map((order) =>
                          React.createElement(
                              "tr",
                              {
                               key: order.id,
                                onClick: () => onSelectOrder(order),
                                 style: {
                                   cursor: "pointer",
                                },
                        } ,

                              React.createElement(
                                  "td",
                                  null,
                                  order.customer_name
                              ),

                              React.createElement(
                                  "td",
                                  null,
                                  formatDate(order.due_date)
                              ),

                              React.createElement(
                                  "td",
                                  null,
                                  formatMoney(order.total_cents)
                              ),

                              React.createElement(
                                  "td",
                                  null,
                                  formatMoney(order.amount_paid_cents)
                              ),

                              React.createElement(
                                  "td",
                                  null,
                                  formatMoney(order.amount_due_cents)
                              ),

                              React.createElement(
                                  "td",
                                  null,
                                  order.status
                              )
                          )
                      )
                  )
              )
    );
}

export default OrderList;