import { useEffect, useState } from "react";
import Button from "react-bootstrap/Button";

const App = () => {
  const [name, setName] = useState("");
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState(null);

  async function loadUsers() {
    const response = await fetch(import.meta.env.VITE_API + "/users");
    const data = await response.json();
    setUsers(data.users);
  }

  async function deleteUser(id) {
    const response = await fetch(import.meta.env.VITE_API + "/users/" + id, {
      method: "DELETE",
    });
    const data = await response.json();
    console.log(data);
    loadUsers();
  }

  async function updateUser(id) {
    const response = await fetch(import.meta.env.VITE_API + "/users/" + id, {
      method: "PUT",
      body: JSON.stringify({ name }),
      headers: {
        "Content-Type": "application/json",
      },
    });
    const data = await response.json();
    console.log(data);
    loadUsers();
    setSelectedUser(null);
    setName("");
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (selectedUser) {
      updateUser(selectedUser._id);
    } else {
      const response = await fetch(import.meta.env.VITE_API + "/users", {
        method: "POST",
        body: JSON.stringify({ name }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      console.log(data);
      setName("");
      loadUsers();
    }
  };

  useEffect(() => {
    loadUsers();
  }, []);

  const handleEdit = (user) => {
    setSelectedUser(user);
    setName(user.name);
  };

  const handleCancelEdit = () => {
    setSelectedUser(null);
    setName("");
  };

  return (
    <div className="pt-2">
      <form className="d-flex" onSubmit={handleSubmit}>
        <div className="form-group">
          <input
            type="name"
            className="form-control"
            placeholder="Coloca aqui un nombre"
            onChange={(e) => setName(e.target.value)}
            value={name}
          />
        </div>
        <button type="submit" className="btn btn-success">
          {selectedUser ? "Actualizar" : "Guardar"}
        </button>

        {selectedUser && (
          <button
            type="button"
            className="btn btn-secondary ml-2"
            onClick={handleCancelEdit}
          >
            Cancelar
          </button>
        )}
      </form>

      <table className="table mt-3">
        <thead>
          <tr>
            <th>ID</th>
            <th>Nombres</th>
            <th>Acciones</th>
          </tr>
        </thead>
        <tbody>
          {users.map((user) => (
            <tr key={user._id}>
              <td>{user._id}</td>
              <td>{user.name}</td>
              <td>
                <Button variant="primary" onClick={() => handleEdit(user)}>
                  Editar
                </Button>
                <button
                  type="button"
                  className="btn btn-danger"
                  onClick={() => deleteUser(user._id)}
                >
                  Eliminar
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default App;
