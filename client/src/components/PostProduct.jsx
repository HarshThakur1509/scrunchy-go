import { useState } from "react";
import { useFormik } from "formik";
import * as Yup from "yup";
import axios from "axios";

export const PostProduct = () => {
  const [message, setMessage] = useState("");

  const formik = useFormik({
    initialValues: {
      name: "",
      price: "",
      image: null,
    },
    validationSchema: Yup.object({
      name: Yup.string().required("Name is required"),
      price: Yup.number()
        .positive("Price must be positive")
        .required("Price is required"),
      image: Yup.mixed().required("Image is required"),
    }),
    onSubmit: async (values) => {
      const formData = new FormData();
      formData.append("name", values.name);
      formData.append("price", values.price);
      formData.append("image", values.image);

      try {
        const response = await axios.post(
          "https://scrunchy.harshthakur.site/api/admin/products",
          formData,
          { withCredentials: true },
          {
            headers: {
              "Content-Type": "multipart/form-data",
            },
          }
        );
        setMessage(response.data.message);
      } catch (error) {
        console.error(error);
        setMessage("Error uploading data");
      }
    },
  });

  const handleFileChange = (event) => {
    formik.setFieldValue("image", event.currentTarget.files[0]);
  };

  return (
    <div className="admin-container">
      <div className="admin-card">
        <h1 className="admin-title">Upload New Product</h1>

        <form className="admin-form" onSubmit={formik.handleSubmit}>
          <div className="input-group">
            <label className="input-label">Product Name</label>
            <input
              type="text"
              name="name"
              className="input-field"
              onChange={formik.handleChange}
              value={formik.values.name}
            />
            {formik.errors.name && (
              <div className="error-message">{formik.errors.name}</div>
            )}
          </div>

          <div className="input-group">
            <label className="input-label">Price</label>
            <input
              type="number"
              name="price"
              className="input-field"
              onChange={formik.handleChange}
              value={formik.values.price}
            />
            {formik.errors.price && (
              <div className="error-message">{formik.errors.price}</div>
            )}
          </div>

          <div className="input-group">
            <label className="input-label">Product Image</label>
            <input
              type="file"
              name="image"
              className="input-field file-input"
              onChange={handleFileChange}
            />
            {formik.errors.image && (
              <div className="error-message">{formik.errors.image}</div>
            )}
          </div>

          <button type="submit" className="admin-button">
            Upload Product
          </button>
        </form>

        {message && <div className="status-message">{message}</div>}
      </div>
    </div>
  );
};
