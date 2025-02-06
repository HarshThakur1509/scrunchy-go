import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { useMutation, useQuery } from "@tanstack/react-query";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import axios from "axios";

export const Home = () => {
  const {
    data: products,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["products"],
    queryFn: async () => {
      const res = await axios.get(
        `https://scrunchy.harshthakur.site/api/products`,
        {
          withCredentials: true,
        }
      );
      return res.data;
    },
    staleTime: 60 * 1000, // Cache for 1 minute
    refetchOnWindowFocus: true,
    retry: 1,
  });

  // Mutation for adding to cart
  const addToCartMutation = useMutation({
    mutationFn: async (id) => {
      await axios.post(
        `https://scrunchy.harshthakur.site/api/cart/add/${id}`,
        null,
        {
          withCredentials: true,
        }
      );
    },
  });

  const handleAddToCart = (id) => {
    addToCartMutation.mutate(id);
  };

  if (isLoading) return <p>Loading videos...</p>;
  if (isError) return <p>Error loading videos: {error.message}</p>;

  return (
    <section className="product-section">
      <div className="container">
        <h2>Featured Products</h2>
        <div className="row">
          {products.map((product) => (
            <div className="card" key={product.ID}>
              <img
                className="card-img-top"
                src={`https://scrunchy.harshthakur.site/api/${product.Image}`}
                alt={product.Name}
              />
              <div className="card-body">
                <h5 className="card-title pro-name">{product.Name}</h5>
                <p className="card-text pro-price">Rs. {product.Price}</p>
                <button
                  className="add-to-cart-btn"
                  onClick={() => handleAddToCart(product.ID)}
                >
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
