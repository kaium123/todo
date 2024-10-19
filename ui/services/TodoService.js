export async function fetchTodos() {
  console.log("Fetching todos..."); 
  const response = await fetch('/api/v1/todos');
  if (!response.ok) throw new Error(`Failed to fetch todos: ${response.status}`);
  const data = await response.json();
  console.log("Fetched todos:", data.data);
  return data.data;
}

export async function createTodo(todo) {
  console.log("Creating todo:", todo); 
  const response = await fetch('/api/v1/todos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(todo),
  });
  if (!response.ok) throw new Error(`Failed to create todo: ${response.status}`);
  const result = await response.json();
  console.log("Created todo response:", result);
  return result;
}

export async function updateTodo(id, todo) {
  console.log(`Updating todo with ID ${id}:`, todo); 
  const response = await fetch(`/api/v1/todos/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(todo),
  });
  if (!response.ok) throw new Error(`Failed to update todo: ${response.status}`);
  const result = await response.json();
  console.log("Updated todo response:", result); 
  return result;
}

export async function deleteTodo(id) {
  console.log(`Attempting to delete todo with ID: ${id}`); 
  const response = await fetch(`/api/v1/todos/${id}`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
  });
  if (!response.ok) {
      console.log(`Failed to delete todo, status: ${response.status}`); 
      throw new Error(`Failed to delete todo: ${response.status}`);
  }
  const result = await response.json();
  console.log("Deleted todo response:", result); 
  return result;
}
