import s1 from "../static/images/s1.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlus, faMinus } from "@fortawesome/free-solid-svg-icons";
import axios from "axios";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const CartItems = ({ cartItems }) => {
  const queryClient = useQueryClient();

  // Mutation for increasing the quantity
  const increaseQuantity = useMutation({
    mutationFn: async ({ id, quantity }) => {
      return await axios.post(
        `http://localhost:3000/cart/quantity/${id}`,
        { quantity },
        { withCredentials: true }
      );
    },
    onSuccess: () => {
      // Invalidate the cart query to refetch updated data
      queryClient.invalidateQueries(["cartItems"]);
    },
  });

  const handleAdd = (id, quantity) => {
    // Mutate the quantity by adding 1
    increaseQuantity.mutate({ id, quantity: quantity + 1 });
  };

  const handleMinus = (id, quantity) => {
    // Mutate the quantity by subtracting 1
    increaseQuantity.mutate({ id, quantity: quantity - 1 });
  };

  const sortedCartItems = [...cartItems].sort((a, b) => a.ID - b.ID);

  return (
    <div className="cart-items">
      {sortedCartItems.map((cartItem) => {
        return (
          <div className="cart-item" key={cartItem.ID}>
            <div className="item-details">
              <img
                className="card-img-top"
                src={`http://localhost:3000/${cartItem.Product.Image}`}
                alt={cartItem.Product.Name}
              />
              <div className="item-info">
                <h5 className="pro-name">{cartItem.Product.Name}</h5>
                <p className="pro-price">
                  Rs.{cartItem.Quantity * cartItem.Product.Price}
                </p>
              </div>
            </div>
            <div className="quantity ml-auto">
              <button
                className="quantity-btn minus"
                onClick={() =>
                  handleMinus(cartItem.Product.ID, cartItem.Quantity)
                }
              >
                <FontAwesomeIcon icon={faMinus} />
              </button>
              <input
                type="number"
                className="quantity-input"
                value={cartItem.Quantity}
                min="1"
                readOnly
              />
              <button
                className="quantity-btn plus"
                onClick={() =>
                  handleAdd(cartItem.Product.ID, cartItem.Quantity)
                }
              >
                <FontAwesomeIcon icon={faPlus} />
              </button>
            </div>
          </div>
        );
      })}
    </div>
  );
};
