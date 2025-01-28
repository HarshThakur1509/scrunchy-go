import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useMutation, useQuery } from "@tanstack/react-query";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import axios from "axios";

export const Home = () => {
  // Fetch products
  const fetchProducts = () => {
    return fetch("http://localhost:3000/products", {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    }).then((res) => res.json());
  };

  const { data, isLoading, isError } = useQuery({
    queryKey: ["products"],
    queryFn: fetchProducts,
  });

  // Mutation for adding to cart
  const addToCartMutation = useMutation({
    mutationFn: async (id) => {
      await axios.post(`http://localhost:3000/cart/add/${id}`, null, {
        withCredentials: true,
      });
    },
  });

  const handleAddToCart = (id) => {
    addToCartMutation.mutate(id);
  };

  if (isLoading) {
    return <p>Loading...</p>;
  }
  if (isError) {
    return <p>Error!</p>;
  }

  return (
    <section className="product-section">
      <div className="container">
        <h2>Featured Products</h2>
        <div className="row">
          {data.map((product) => (
            <div className="card" key={product.ID}>
              <img
                className="card-img-top"
                src={`http://localhost:3000/${product.Image}`}
                alt={product.Name}
              />
              <div className="card-body">
                <h5 className="card-title pro-name">{product.Name}</h5>
                <p className="card-text pro-price">Rs. {product.Price}</p>
                <button onClick={() => handleAddToCart(product.ID)}>
                  <FontAwesomeIcon icon={faPlus} />
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};
