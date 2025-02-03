import { CartItems } from "../components/CartItems";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import axios from "axios";

export const Cart = () => {
  // Fetch cart items from the API
  const {
    data: cart,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["cart"],
    queryFn: async () => {
      const response = await axios.get(
        "https://scrunchy.harshthakur.site/api/cart",
        {
          withCredentials: true,
        }
      );

      return response.data;
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error fetching cart items...</div>;

  // Calculate the total amount of the cart
  // const TotalAmount = (items) => {
  //   let total = 0;
  //   items.forEach((item) => {
  //     total += item.Quantity * item.Product.Price;
  //   });
  //   return total;
  // };

  return (
    <section className="Cart">
      <div className="container">
        <h2>Your Cart</h2>
        <CartItems cartItems={cart.CartItems} />
        <div className="cart-total">
          <h4>Total: Rs. {cart.Total}</h4>
          <button className="btn btn-primary">
            <Link to="/checkout" state={{ cart }}>
              Checkout
            </Link>
          </button>
        </div>
      </div>
    </section>
  );
};
