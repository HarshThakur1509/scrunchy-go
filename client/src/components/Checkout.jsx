import { useLocation } from "react-router-dom";
import { useContext } from "react";
import { useForm } from "react-hook-form";
import * as yup from "yup";
import { yupResolver } from "@hookform/resolvers/yup";
import { LoginContext } from "../App";
import axios from "axios";

export const Checkout = () => {
  const location = useLocation();
  const { cart } = location.state || {};
  const { userDetails } = useContext(LoginContext);

  const schema = yup.object().shape({
    city: yup.string().min(4).max(20).required("City Required"),
    location: yup.string().min(4).max(30).required("Location Required"),
    state: yup.string().min(4).max(20).required("State Required"),
    zipCode: yup.string().min(4).max(20).required("ZipCode Required"),
  });

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm({
    resolver: yupResolver(schema),
  });

  const onSubmit = async (formdata) => {
    try {
      await axios.put(
        "https://scrunchy.harshthakur.site/api/address",
        formdata,
        {
          withCredentials: true,
        }
      );
      window.location.reload();
    } catch (err) {
      console.log(err);
    }
  };

  const handlePayment = async () => {
    // Fetch order details from the server
    try {
      const response = await fetch(
        "https://scrunchy.harshthakur.site/api/pay",
        {
          method: "GET",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
        }
      );

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
              "https://scrunchy.harshthakur.site/api/payresponse",
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
    <section className="checkout-section">
      <div className="container">
        <h2>Checkout</h2>
        <div className="checkout-content">
          {/* Shipping and Payment Form */}
          <div className="checkout-form">
            <h3>Shipping Details</h3>
            <form className="address-form" onSubmit={handleSubmit(onSubmit)}>
              <div className="input-group">
                <input
                  type="text"
                  id="location"
                  defaultValue={userDetails?.Address?.Location}
                  placeholder="Location"
                  {...register("location")}
                  className={`input-field ${
                    errors.location ? "input-error" : ""
                  }`}
                />
              </div>
              <div className="input-group">
                <input
                  type="text"
                  id="city"
                  defaultValue={userDetails?.Address?.City}
                  placeholder="City"
                  {...register("city")}
                  className={`input-field ${errors.city ? "input-error" : ""}`}
                />
              </div>
              <div className="input-group">
                <input
                  type="text"
                  id="state"
                  defaultValue={userDetails?.Address?.State}
                  placeholder="State"
                  {...register("state")}
                  className={`input-field ${errors.state ? "input-error" : ""}`}
                />
              </div>
              <div className="input-group">
                <input
                  type="text"
                  id="zipCode"
                  defaultValue={userDetails?.Address?.ZipCode}
                  placeholder="ZIP code"
                  {...register("zipCode")}
                  className={`input-field ${
                    errors.zipCode ? "input-error" : ""
                  }`}
                />
              </div>
              <button
                className="address-button"
                type="submit"
                disabled={isSubmitting}
              >
                {isSubmitting ? "Submitting..." : "Submit"}
              </button>
            </form>
          </div>

          {/* Order Summary */}
          <div className="order-summary">
            <h3>Order Summary</h3>
            <div className="cart-items">
              {cart.CartItems.map((item) => (
                <div className="cart-item" key={item.ID}>
                  <img
                    src={`https://scrunchy.harshthakur.site/api/${item.Product.Image}`}
                    alt={item.Product.Name}
                    className="item-image"
                  />
                  <div className="item-info">
                    <p className="pro-name">{item.Product.Name}</p>
                    <p className="pro-price">Rs. {item.Price}</p>
                    <p>Quantity: {item.Quantity}</p>
                  </div>
                </div>
              ))}
            </div>
            <div className="cart-total">
              <p>Total:</p>
              <p className="total-price">Rs. {cart.Total}</p>
            </div>
            <button
              type="submit"
              className="checkout-btn"
              onClick={handlePayment}
            >
              Place Order
            </button>
          </div>
        </div>
      </div>
    </section>
  );
};
