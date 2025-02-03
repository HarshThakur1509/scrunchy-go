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
    <div>
      <h1>Upload Form</h1>
      <form onSubmit={formik.handleSubmit}>
        <div>
          <label>Name</label>
          <input
            type="text"
            name="name"
            onChange={formik.handleChange}
            value={formik.values.name}
          />
          {formik.errors.name && <div>{formik.errors.name}</div>}
        </div>
        <div>
          <label>Price</label>
          <input
            type="number"
            name="price"
            onChange={formik.handleChange}
            value={formik.values.price}
          />
          {formik.errors.price && <div>{formik.errors.price}</div>}
        </div>
        <div>
          <label>Image</label>
          <input type="file" name="image" onChange={handleFileChange} />
          {formik.errors.image && <div>{formik.errors.image}</div>}
        </div>
        <button type="submit">Submit</button>
      </form>
      {message && <div>{message}</div>}
    </div>
  );
};
