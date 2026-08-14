import React, { useEffect, useState } from "react";
import Login from "./components/Login";
import OrderList from "./components/OrderList";
import OrderDetails from "./components/OrderDetails";
import { getCurrentUser, createOrder} from "./api";
import OrderForm from "./components/OrderForm";

function App() {
    const [authenticated, setAuthenticated] = useState(false);
    const [loading, setLoading] = useState(true);
    const [selectedOrderId, setSelectedOrderId] = useState(null);
    const [creatingOrder, setCreatingOrder] = useState(false);

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
        return React.createElement("p", null, "Loading...");
    }

    if (!authenticated) {
        return React.createElement(Login, {
            onLogin: () => setAuthenticated(true),
        });
    }

    if (selectedOrderId) {
        return React.createElement(OrderDetails, {
            orderId: selectedOrderId,
            onBack: () => setSelectedOrderId(null),
        });
    }

    if (creatingOrder) {
    return React.createElement(OrderForm, {
        onSubmit: async (order) => {
            try {
                await createOrder(order);
                setCreatingOrder(false);
            } catch (error) {
                console.error(error);
            }
        },

        onCancel: () => {
            setCreatingOrder(false);
        },
    });
}

    return React.createElement(OrderList, {
    onSelectOrder: (order) => {
        setSelectedOrderId(order.id);
    },

    onCreateOrder: () => {
        setCreatingOrder(true);
    },
});
}

export default App;