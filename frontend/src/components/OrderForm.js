import React, { useState } from "react";

function OrderForm({ onSubmit, onCancel }) {
    const [customerName, setCustomerName] = useState("");
    const [dueDate, setDueDate] = useState("");

    const [items, setItems] = useState([
        {
            description: "",
            quantity: 1,
            unit_price_cents: 0,
        },
    ]);

    function updateItem(index, field, value) {
        setItems((current) =>
            current.map((item, i) => {
                if (i !== index) {
                    return item;
                }

                return {
                    ...item,
                    [field]:
                        field === "description"
                            ? value
                            : Number(value),
                };
            })
        );
    }

    function addItem() {
        setItems((current) => [
            ...current,
            {
                description: "",
                quantity: 1,
                unit_price_cents: 0,
            },
        ]);
    }

    function removeItem(index) {
        setItems((current) =>
            current.filter((_, i) => i !== index)
        );
    }

    function handleSubmit(event) {
        event.preventDefault();

        onSubmit({
            customer_name: customerName,
            due_date: new Date(dueDate).toISOString(),
            line_items: items,
        });
    }

    return React.createElement(
        "div",
        { style: { padding: "40px" } },

        React.createElement(
            "h1",
            null,
            "Create Order"
        ),

        React.createElement(
            "form",
            { onSubmit: handleSubmit },

            React.createElement(
                "div",
                null,

                React.createElement(
                    "label",
                    null,
                    "Customer Name"
                ),

                React.createElement("input", {
                    type: "text",
                    value: customerName,
                    onChange: (event) =>
                        setCustomerName(event.target.value),
                    required: true,
                })
            ),

            React.createElement(
                "div",
                null,

                React.createElement(
                    "label",
                    null,
                    "Due Date"
                ),

                React.createElement("input", {
                    type: "datetime-local",
                    value: dueDate,
                    onChange: (event) =>
                        setDueDate(event.target.value),
                    required: true,
                })
            ),

            React.createElement(
                "h2",
                null,
                "Line Items"
            ),

            items.map((item, index) =>
                React.createElement(
                    "div",
                    {
                        key: index,
                        style: {
                            marginBottom: "20px",
                            padding: "12px",
                            border: "1px solid #ddd",
                        },
                    },

                    React.createElement(
                        "div",
                        null,

                        React.createElement(
                            "label",
                            null,
                            "Description"
                        ),

                        React.createElement("input", {
                            type: "text",
                            value: item.description,
                            onChange: (event) =>
                                updateItem(
                                    index,
                                    "description",
                                    event.target.value
                                ),
                            required: true,
                        })
                    ),

                    React.createElement(
                        "div",
                        null,

                        React.createElement(
                            "label",
                            null,
                            "Quantity"
                        ),

                        React.createElement("input", {
                            type: "number",
                            min: "1",
                            value: item.quantity,
                            onChange: (event) =>
                                updateItem(
                                    index,
                                    "quantity",
                                    event.target.value
                                ),
                            required: true,
                        })
                    ),

                    React.createElement(
                        "div",
                        null,

                        React.createElement(
                            "label",
                            null,
                            "Unit Price (cents)"
                        ),

                        React.createElement("input", {
                            type: "number",
                            min: "0",
                            value: item.unit_price_cents,
                            onChange: (event) =>
                                updateItem(
                                    index,
                                    "unit_price_cents",
                                    event.target.value
                                ),
                            required: true,
                        })
                    ),

                    items.length > 1 &&
                        React.createElement(
                            "button",
                            {
                                type: "button",
                                onClick: () =>
                                    removeItem(index),
                            },
                            "Remove"
                        )
                )
            ),

            React.createElement(
                "button",
                {
                    type: "button",
                    onClick: addItem,
                },
                "+ Add Item"
            ),

            React.createElement(
                "div",
                {
                    style: {
                        marginTop: "24px",
                    },
                },

                React.createElement(
                    "button",
                    {
                        type: "submit",
                    },
                    "Create Order"
                ),

                React.createElement(
                    "button",
                    {
                        type: "button",
                        onClick: onCancel,
                        style: {
                            marginLeft: "10px",
                        },
                    },
                    "Cancel"
                )
            )
        )
    );
}

export default OrderForm;