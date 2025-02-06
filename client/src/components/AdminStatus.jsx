import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useState, useEffect, useContext } from "react";
import { LoginContext } from "../App";

export const AdminStatus = () => {
  const { userDetails } = useContext(LoginContext);
  const [list, setList] = useState([]);

  // Fetch users list
  const users = async () => {
    const response = await axios.get(
      "https://scrunchy.harshthakur.site/api/admin/listusers",
      {
        withCredentials: true,
      }
    );
    console.log(response.data);

    return response.data;
  };

  const { data, isLoading, isError } = useQuery({
    queryKey: ["users"],
    queryFn: users,
    onSuccess: (fetchedData) => {
      // Directly process data after fetching
      const filteredData = fetchedData.filter(
        (user) => user.ID !== userDetails.ID // Exclude current user
      );
      setList(filteredData);
    },
  });

  // Filter to keep only specified attributes (ID, Email, Admin)
  const keepAttributes = (data, attributes) => {
    return data.map((obj) =>
      attributes.reduce((acc, key) => {
        if (key in obj) acc[key] = obj[key];
        return acc;
      }, {})
    );
  };

  useEffect(() => {
    if (data && data.length > 0) {
      // Apply attribute filtering and remove current user
      const filteredData = keepAttributes(
        data.filter((user) => user.ID !== userDetails.ID),
        ["ID", "Email", "Admin"]
      );
      setList(filteredData);
    }
  }, [data, userDetails]);

  // Handle checkbox change for admin status
  const handleCheckboxChange = (id, isAdmin) => {
    setList((prevList) =>
      prevList.map((user) =>
        user.ID === id ? { ...user, Admin: isAdmin } : user
      )
    );
  };

  const handleSetButtonClick = async (id, isAdmin) => {
    try {
      await axios.post(
        `https://scrunchy.harshthakur.site/api/admin/status/${id}`,
        { admin: isAdmin }, // Include the updated admin status
        { withCredentials: true }
      );
    } catch (error) {
      console.error(`Error updating user ${id}:`, error);
    }
  };

  const handleDeleteUser = async (id) => {
    try {
      await axios.delete(
        `https://scrunchy.harshthakur.site/api/admin/delete/${id}`,
        {
          withCredentials: true,
        }
      );

      // Update the state to remove the deleted user
      setList((prevList) => prevList.filter((user) => user.ID !== id));
    } catch (error) {
      console.error(`Error deleting user ${id}:`, error);
    }
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (isError) {
    return <div>Error fetching data...</div>;
  }

  const listUsers = list.map((user) => (
    <div key={user.ID}>
      <p>Email: {user.Email}</p>
      <input
        type="checkbox"
        checked={user.Admin}
        onChange={(e) => handleCheckboxChange(user.ID, e.target.checked)}
      />
      <button onClick={() => handleSetButtonClick(user.ID, user.Admin)}>
        Set
      </button>
      <button onClick={() => handleDeleteUser(user.ID)}>Delete</button>
    </div>
  ));

  return <div>{listUsers}</div>;
};
