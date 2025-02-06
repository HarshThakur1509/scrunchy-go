import { Link } from "react-router-dom";

export const Admin = () => {
  return (
    <div className="Admin">
      <Link to="/admin/user">User Status</Link>
      <Link to="/admin/product">Post Product</Link>
    </div>
  );
};
