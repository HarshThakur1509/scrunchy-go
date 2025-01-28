import { useLocation } from "react-router-dom";
import { useContext } from "react";
import { LoginContext } from "../App";

export const Checkout = () => {
  const location = useLocation();
  const { cart } = location.state || {};
  const { userDetails } = useContext(LoginContext);

  const handlePayment = async () => {
    // Fetch order details from the server
    try {
      const response = await fetch("http://localhost:3000/pay", {
        method: "GET",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch payment details");
      }

      const payData = await response.json();

      // Razorpay options
      const options = {
        key: "rzp_test_MQUwdShJLMIpOu", // Enter the Key ID generated from the Razorpay Dashboard
        amount: cart.Total * 100, // Amount is in currency subunits
        currency: "INR",
        name: "Scrunchy",
        description: "Scrunchy Payment",
        image: "https://example.com/your_logo",
        order_id: payData.id, // Set the order ID received from the server
        handler: async (response) => {
          // Handle the successful payment response
          const paymentData = {
            payment_id: response.razorpay_payment_id,
            order_id: response.razorpay_order_id,
            signature: response.razorpay_signature,
          };

          try {
            const paymentResponse = await fetch(
              "http://localhost:3000/payresponse",
              {
                method: "POST",
                credentials: "include",
                headers: {
                  "Content-Type": "application/json",
                },
                body: JSON.stringify(paymentData),
              }
            );

            if (!paymentResponse.ok) {
              throw new Error("Failed to process payment");
            }

            alert("Payment successfully!");
          } catch (error) {
            alert("Payment verification failed: " + error.message);
          }
        },
        prefill: {
          name: userDetails.Name,
          email: userDetails.Email,
          contact: userDetails.Phone,
        },
        notes: {
          address: "Razorpay Corporate Office",
        },
        theme: {
          color: "#3399cc",
        },
      };

      // Initialize Razorpay
      const rzp = new window.Razorpay(options);

      rzp.on("payment.failed", (response) => {
        alert(`Payment failed: ${response.error.description}`);
      });

      rzp.open();
    } catch (error) {
      alert("Failed to initialize payment: " + error.message);
    }
  };

  return (
    <div>
      Checkout
      {/* {cart ? (
        <div>
          <p>ID: {cart.ID}</p>
        </div>
      ) : (
        <p>No user data available.</p>
      )} */}
      <button onClick={handlePayment}>Order</button>
    </div>
  );
};
